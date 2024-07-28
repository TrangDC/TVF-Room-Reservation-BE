package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController interface {
	OAuthLogin(c *gin.Context)
	OAuthCallback(c *gin.Context)
	RefreshToken(c *gin.Context)
	OAuthLogout(c *gin.Context)
	OAuthValidate(c *gin.Context)
}

// authCtrlImpl is the implementation of AuthController.
type authCtrlImpl struct {
	authService       service.AuthService
	logger            *zap.Logger
	clientRedirectURL string
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService service.AuthService, clientRedirectURL string, logger *zap.Logger) AuthController {
	return &authCtrlImpl{
		authService:       authService,
		clientRedirectURL: clientRedirectURL,
		logger:            logger,
	}
}

// OAuthLogin handles the OAuth login.
func (i authCtrlImpl) OAuthLogin(c *gin.Context) {
	redirectUrl, err := i.authService.GetRedirectLoginEndpoint(c, c.Request, c.Writer)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, redirectUrl)
}

// OAuthCallback handles the OAuth callback.
func (i authCtrlImpl) OAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()
	if err := c.Query("error"); len(err) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	r := service.OAuthLoginCallbackRequest{}
	if err := c.ShouldBindQuery(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.WrapGQLBadRequestError(ctx, "Invalid request to create token"))
		return
	}
	token, err := i.authService.GetToken(ctx, r, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.WrapGQLBadRequestError(ctx, "Invalid request to create token"))
		return
	}
	queryParams := url.Values{
		"accessToken":  {token.AccessToken},
		"refreshToken": {token.RefreshToken},
		"expiresAt":    {fmt.Sprintf("%d", token.ExpiresAt.Unix())},
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?%s", i.clientRedirectURL, queryParams.Encode()))
}

// RefreshToken handles the refresh token.
func (i authCtrlImpl) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	r := service.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.WrapGQLBadRequestError(ctx, "Invalid request to refresh token with invalid body"))
		return
	}

	token, err := i.authService.RefreshToken(ctx, r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, util.WrapGQLBadRequestError(ctx, "Invalid request to refresh token"))
		return
	}

	c.JSON(http.StatusOK, token)
}

func (i authCtrlImpl) OAuthLogout(c *gin.Context) {
	redirectUrl, err := i.authService.Logout(c, c.Request, c.Writer)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, redirectUrl)
}

func (i authCtrlImpl) OAuthValidate(c *gin.Context) {
	ctx := c.Request.Context()
    ok, err := i.authService.ValidateToken(ctx, c.Request)
    if err!= nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, util.WrapGQLUnauthorizedError(ctx))
        return
    }

    c.JSON(http.StatusOK, ok)
}
