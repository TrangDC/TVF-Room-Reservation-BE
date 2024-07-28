package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/config"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/azuread"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/pg"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/middleware"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/resolver"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/rest"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/service"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewServerCmd(configs *config.Configurations, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "run api server",
		Long:  "run api server with graphql",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			defer func() {
				err := recover()
				if err != nil {
					logger.Fatal("recover error", zap.Any("error", err))
				}
			}()

			// Create PostgreSQL connection
			db, err := pg.NewDBConnection(configs.Postgres, logger)
			if err != nil {
				logger.Error("Getting error connect to postgresql database", zap.Error(err))
				os.Exit(1)
			}
			defer db.Close()
			// Create ent client
			entDriver := entsql.OpenDB("postgres", db)

			entOptions := []ent.Option{
				ent.Driver(entDriver),
			}
			if configs.App.Debug {
				entOptions = append(entOptions, ent.Debug())
			}

			entClient := ent.NewClient(entOptions...)

			defer entClient.Close()

			// Create validator
			validator := validator.New()
			// Add translator for validator
			en := en.New()
			uni := ut.New(en, en)
			validationTranslator, _ := uni.GetTranslator("en")
			// Register default translation for validator
			err = en_translations.RegisterDefaultTranslations(validator, validationTranslator)
			if err != nil {
				logger.Error("Getting error from register default translation", zap.Error(err))
				os.Exit(1)
			}

			// Create sessions store with cookie
			sessionStore := sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))
			sessionStore.MaxAge(60)

			// Create AzureAD OAuth client
			var azureADOAuthClient azuread.AzureADOAuth
			// Authentication with AzureAD
			if configs.AzureADOAuth.Enabled {
				azureADOAuthClient, err = azuread.NewAzureADOAuth(configs.AzureADOAuth, sessionStore)
				if err != nil {
					logger.Error("Getting error create to AzureAD oauth", zap.Error(err))
					os.Exit(1)
				}
			}
			serviceRegistry := service.NewService(azureADOAuthClient, entClient, logger)
			restController := rest.NewRestController(serviceRegistry, configs.AzureADOAuth.ClientRedirectUrl, logger)

			// GraphQL schema resolver handler.
			resolverHandler := handler.NewDefaultServer(resolver.NewSchema(serviceRegistry, entClient, validator, validationTranslator, logger))

			// Handler for GraphQL Playground
			playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")

			if !configs.App.Debug {
				gin.SetMode(gin.ReleaseMode)
			}
			r := gin.New()
			// Handle a method not allowed.
			r.HandleMethodNotAllowed = true

			// Use middlewares
			r.Use(
				ginzap.Ginzap(logger, time.RFC3339, true),
				ginzap.RecoveryWithZap(logger, true),
				middleware.CorsMiddleware(),
				middleware.RequestCtxMiddleware(),
			)

			readyRouter := r.Group("ready")
			{
				readyRouter.GET("/readiness", middleware.ReadinessCheckMiddleware(db, logger))
				readyRouter.GET("/liveliness", middleware.LivelinessCheckMiddleware(db, logger))
			}

			if configs.AzureADOAuth.Enabled {
				authRouter := r.Group("/auth")
				{
					authRouter.GET("/login", restController.Auth().OAuthLogin)
					authRouter.GET("/callback", restController.Auth().OAuthCallback)
					authRouter.POST("/refresh", restController.Auth().RefreshToken)
					authRouter.GET("/logout", restController.Auth().OAuthLogout)
					authRouter.POST("/readiness", restController.Auth().OAuthValidate)
				}
			}

			graphqlRouter := r.Group("/graphql")
			{
				if configs.AzureADOAuth.Enabled {
					graphqlRouter.Use(middleware.AuthenticateMiddleware(azureADOAuthClient, db, logger))
				}

				graphqlRouter.POST("", func(c *gin.Context) {
					resolverHandler.ServeHTTP(c.Writer, c.Request)
				})

				graphqlRouter.OPTIONS("", func(c *gin.Context) {
					c.Status(200)
				})
			}

			if configs.App.Debug {
				r.GET("/", func(c *gin.Context) {
					playgroundHandler.ServeHTTP(c.Writer, c.Request)
				})
			}

			server := &http.Server{
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 30 * time.Second,
				Addr:         ":8080",
				Handler:      r,
			}

			// Graceful shutdown
			idleConnectionsClosed := make(chan struct{})
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

				<-c

				ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()
				// A interrupt signal has sent to us, let's shutdown server with gracefully
				logger.Debug("Stopping server...")

				if err := server.Shutdown(ctx); err != nil {
					logger.Error("Graceful shutdown has failed with error: %s", zap.Error(err))
				}

				if err := db.Close(); err != nil {
					logger.Error("Closing db connection has error", zap.Error(err))
				}

				close(idleConnectionsClosed)
			}()

			go func() {
				logger.Debug("Listing on the port: 8080")
				if err := server.ListenAndServe(); err != http.ErrServerClosed {
					logger.Error("Run server has error", zap.Error(err))
					// Exit the application if run fail
					os.Exit(1)
				} else {
					logger.Info("Server was closed by shutdown gracefully")
				}
			}()

			<-idleConnectionsClosed
		},
	}
}
