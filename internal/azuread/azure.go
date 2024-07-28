package azuread

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/config"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/dto"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/models"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	OAUTH_STATE_NAME = "AD_OAUTH_STATE"
	OAUTH_NONCE_NAME = "AD_OAUTH_NONCE"
)

type AzureADOAuth interface {
	GetOAuthLoginEndpoint(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error)
	GetToken(ctx context.Context, token *oauth2.Token) (JwtToken, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (JwtToken, error)
	GetOAuthTokenFromCallback(ctx context.Context, r *http.Request, code, state string) (*oauth2.Token, error)
	ValidateNonce(ctx context.Context, r *http.Request, nonce string) error
	VerifyAccessToken(ctx context.Context, token string) error
	DecodeToken(ctx context.Context, token string) (*models.Token, error)
	Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error)
}

type JwtToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
	Email        string    `json:"email"`
}

type JwtClaims struct {
	Email string `json:"email"`
}

type azureADOAuthImpl struct {
	provider      *oidc.Provider
	oauthConfig   *oauth2.Config
	oidcConfig    *oidc.Config
	sessionsStore *sessions.CookieStore
}

func NewAzureADOAuth(config config.AzureADOAuthConfig, sessionsStore *sessions.CookieStore) (AzureADOAuth, error) {
	//discovery the oidc endpoint
	if config.TenantID == "" {
		return nil, errors.WithStack(errors.New("tenantID is required"))
	}
	provider, err := oidc.NewProvider(context.Background(), fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", config.TenantID))
	if err != nil {
		return nil, err
	}

	scopes := []string{oidc.ScopeOpenID, "profile", "email", "offline_access"}
	if config.Scopes != "" {
		scopes = append(scopes, config.Scopes)
	}

	return &azureADOAuthImpl{
		oidcConfig: &oidc.Config{
			ClientID: config.ClientID,
		},
		provider: provider,
		oauthConfig: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectUrl,
			// discovery returns the OAuth2 endpoints.
			Endpoint: provider.Endpoint(),
			// "openid" is a required scope for OpenID Connect flows.
			Scopes: scopes,
		},
		sessionsStore: sessionsStore,
	}, nil

}

func (i azureADOAuthImpl) GetOAuthLoginEndpoint(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	state, err := i.setCallBackSession(w, r, OAUTH_STATE_NAME)
	if err != nil {
		return "", errors.WithStack(fmt.Errorf("create oauth state with session has error: %s", err.Error()))
	}

	nonce, err := i.setCallBackSession(w, r, OAUTH_NONCE_NAME)
	if err != nil {
		return "", errors.WithStack(fmt.Errorf("create oauth nonce with session has error: %s", err.Error()))
	}

	return i.oauthConfig.AuthCodeURL(state, oidc.Nonce(nonce)), nil
}

func (i azureADOAuthImpl) VerifyAccessToken(ctx context.Context, token string) error {
	_, err := i.provider.Verifier(i.oidcConfig).Verify(ctx, token)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to verify access token: %s", err.Error()))
	}

	return nil
}

func (i azureADOAuthImpl) GetToken(ctx context.Context, token *oauth2.Token) (JwtToken, string, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return JwtToken{}, "", errors.WithStack(errors.New("getting invalid id_token"))
	}

	idToken, err := i.provider.Verifier(i.oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		return JwtToken{}, "", errors.WithStack(fmt.Errorf("failed to verify id_token: %s", err.Error()))
	}

	claims := JwtClaims{}

	if err := idToken.Claims(&claims); err != nil {
		return JwtToken{}, "", errors.WithStack(fmt.Errorf("failed to parse id_token claims: %s", err.Error()))
	}

	return JwtToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
		Email:        claims.Email,
	}, idToken.Nonce, nil
}

func (i azureADOAuthImpl) RefreshToken(ctx context.Context, refreshToken string) (JwtToken, error) {
	oauthToken, err := i.oauthConfig.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	}).Token()

	if err != nil {
		return JwtToken{}, errors.WithStack(fmt.Errorf("failed to refresh token: %s", err.Error()))
	}

	token, _, err := i.GetToken(ctx, oauthToken)
	if err != nil {
		return JwtToken{}, err
	}

	return token, nil
}

func (i azureADOAuthImpl) GetOAuthTokenFromCallback(ctx context.Context, r *http.Request, code, state string) (*oauth2.Token, error) {
	session, _ := i.sessionsStore.Get(r, "oauth")
	sessionState := session.Values[OAUTH_STATE_NAME]
	if sessionState == "" {
		return nil, errors.WithStack(fmt.Errorf("getting invalid session state"))
	}

	if sessionState != state {
		return nil, errors.WithStack(fmt.Errorf("getting invalid session state: %s", state))
	}

	token, err := i.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to exchange token: %s", err.Error()))
	}

	return token, nil
}

func (i azureADOAuthImpl) ValidateNonce(ctx context.Context, r *http.Request, nonce string) error {
	session, _ := i.sessionsStore.Get(r, "oauth")
	sessionNonce := session.Values[OAUTH_NONCE_NAME]
	if sessionNonce == "" {
		return errors.WithStack(fmt.Errorf("getting invalid session nonce"))
	}

	if sessionNonce != nonce {
		return errors.WithStack(fmt.Errorf("getting invalid nonce from session: %s", nonce))
	}

	return nil
}

func (i azureADOAuthImpl) setCallBackSession(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	value := uuid.New().String()

	session, _ := i.sessionsStore.Get(r, "oauth")
	session.Values[name] = value

	if err := session.Save(r, w); err != nil {
		return "", errors.WithStack(fmt.Errorf("save session has error: %s", err.Error()))
	}

	return value, nil
}

func (i azureADOAuthImpl) DecodeToken(ctx context.Context, tokenString string) (*models.Token, error) {
	tokenAu, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenAu.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	tokenData := dto.ConvertTokenClaim(claims)
	return tokenData, err
}

func (i azureADOAuthImpl) Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
	session, _ := i.sessionsStore.Get(r, "oauth")
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		return "", errors.WithStack(fmt.Errorf("save session has error: %s", err.Error()))
	}

	azureTenantID := "8f413df9-2f7e-4d51-9384-9cea3a8d5874"
	postLogoutRedirectURI := "https://grateful-yeti-deep.ngrok-free.app/auth/login"
	encodedRedirectURI := url.QueryEscape(postLogoutRedirectURI)
	redirectUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/logout?post_logout_redirect_uri=%s", azureTenantID, encodedRedirectURI)

	return redirectUrl, nil
}
