package dto

import (
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/models"

	"github.com/golang-jwt/jwt"
)

func ConvertTokenClaim(claims jwt.MapClaims) *models.Token {
	return &models.Token{
		Audience:           claims["aud"].(string),
		Issuer:             claims["iss"].(string),
		IssuedAt:           int64(claims["iat"].(float64)),
		NotBefore:          int64(claims["nbf"].(float64)),
		Expiration:         int64(claims["exp"].(float64)),
		AuthenticationInfo: claims["aio"].(string),
		AuthorizedParty:    claims["azp"].(string),
		AuthorizedPartyCR:  claims["azpacr"].(string),
		Name:               claims["name"].(string),
		ObjectID:           claims["oid"].(string),
		PreferredUsername:  claims["preferred_username"].(string),
		RefreshTokenHash:   claims["rh"].(string),
		Scopes:             claims["scp"].(string),
		Subject:            claims["sub"].(string),
		TenantID:           claims["tid"].(string),
		UniqueTokenID:      claims["uti"].(string),
		Version:            claims["ver"].(string),
	}
}
