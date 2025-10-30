package introspect

import (
	"os"
	"time"

	"local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/models"
	utils "local/bomboclat-oauth-server/utils"

	"github.com/golang-jwt/jwt/v5"
)

func (as *IntrospectService) Introspect(m *models.InstrospectModelInput) (*map[string]any, error) {

	switch m.TokenTypeHint {
	case "access_token":
		tokenHash := utils.HashToken256(m.Token)

		tokenData, err := database.IntrospectAccessToken(as.DBConn, tokenHash)

		if err != nil {
			return nil, err
		}

		if tokenData == nil {
			return nil, &utils.NoAccessTokenFoundError{}
		}

		token, err := jwt.ParseWithClaims(m.Token, &models.CustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return nil, &utils.TokenParsingError{}
		}

		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {

			isTokenActive := !tokenData.Revoked && claims.ExpiresAt > time.Now().UTC().Unix()

			m := map[string]any{
				"client_id":  tokenData.ClientId,
				"scope":      tokenData.Scopes,
				"active":     isTokenActive,
				"sub":        claims.Subject,
				"exp":        claims.ExpiresAt,
				"iat":        claims.IssuedAt,
				"iss":        claims.Issuer,
				"token_type": m.TokenTypeHint,
				"aud":        claims.Audience,
			}

			return &m, nil
		}

		return nil, &utils.TokenParsingError{}
	case "refresh_token":
		tokenHash := utils.HashToken256(m.Token)

		tokenData, err := database.IntrospectRefreshToken(as.DBConn, tokenHash)

		if err != nil {
			return nil, err
		}

		if tokenData == nil {
			return nil, &utils.NoAccessTokenFoundError{}
		}

		isTokenActive := !tokenData.Revoked && tokenData.ExpiresAt.Unix() > time.Now().UTC().Unix()
		exp := tokenData.ExpiresAt.Unix()
		iat := tokenData.CreatedAt.Unix()

		m := map[string]any{
			"client_id":  tokenData.ClientId,
			"scope":      tokenData.Scopes,
			"active":     isTokenActive,
			"token_type": m.TokenTypeHint,
			"exp":        exp,
			"iat":        iat,
		}

		return &m, err

	default:
		return nil, nil
	}
}
