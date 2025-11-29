package sessions

import (
	"errors"
	"fmt"
	custom_errors "local/bomboclat-oidc-service/errors"
	custom_types "local/bomboclat-oidc-service/types"
	utils "local/bomboclat-oidc-service/utils"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (ss *SessionService) EndSession(m custom_types.EndSessionInput) error {

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("JWT_RSA_PUBLIC_KEY")))
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(m.IdTokenHint, &custom_types.IDTokenClaims{}, func(token *jwt.Token) (any, error) {
		return key, nil
	})

	if err != nil {
		return custom_errors.TokenParsingError(err)
	}

	if claims, ok := token.Claims.(*custom_types.IDTokenClaims); ok {

		isIssuer := claims.Iss == os.Getenv("OIDC_BASE_URL")

		if !isIssuer {
			return custom_errors.TokenIssuerMismatchError(nil)
		}

		userCookie := m.UserCookie

		//remove redis session
		sessionId := userCookie.Value
		if err := utils.DeleteHashAll(ss.RedisClient, "user_session:"+sessionId); err != nil {
			var appError *custom_errors.AppError
			if errors.As(err, &appError) && appError.Code == "REDIS_DELETE_ERROR" {
				return custom_errors.UserLogoutError(appError.Cause)
			}
		}

		fmt.Print("Success logout")

		return nil
	}

	return custom_errors.TokenParsingError(nil)
}
