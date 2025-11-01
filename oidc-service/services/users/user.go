package users

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"local/bomboclat-oidc-service/database"
	custom_types "local/bomboclat-oidc-service/types"
	utils "local/bomboclat-oidc-service/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/golang-jwt/jwt/v5"
)

func (us *UserService) Login(userDetails custom_types.UserLoginDetails) (*http.Cookie, error) {

	data, err := database.FindUserByEmailAndPasswordHash(us.DBConn, context.Background(),
		userDetails.Email, utils.HashToken256(userDetails.Password))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Invalid user details")
			return nil, &utils.UserNotFoundError{}
		}
		return nil, err
	}

	if data == nil {
		log.Print("No user with this email!")
		return nil, &utils.UserNotFoundError{}
	}

	sessionID := uuid.New().String()

	userDetailsMap := map[string]string{
		"user_id": data.UUID,
		"scope":   "deny",
	}

	if err := utils.SetValueToHash(us.RedisClient, "user_session:"+sessionID, userDetailsMap); err != nil {
		log.Print(err)
		return nil, &utils.RedisSetHasError{}
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func (us *UserService) Register(details custom_types.UserRegistrationDetails) error {
	// TODO: Invalidate if user with this email already exsists
	user, err := database.FindUserByEmail(us.DBConn, context.Background(), details.Email)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if user != nil {
		return &utils.UserAlreadyExistsError{}
	}

	// TODO: Create new user if not
	err = database.InsertUser(us.DBConn, &details)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UserInfo(authToken string) (*custom_types.UserInfo, error) {

	token, err := jwt.ParseWithClaims(authToken, &custom_types.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, &utils.TokenParsingError{}
	}

	if claims, ok := token.Claims.(*custom_types.CustomClaims); ok && token.Valid {

		userUUID := claims.Subject

		user, err := database.FindUserByUUID(us.DBConn, context.Background(), userUUID)

		userInfo := &custom_types.UserInfo{
			Username: user.Username,
			Email:    user.Email,
		}

		if err != nil {
			return nil, err
		}

		return userInfo, nil
	}

	return nil, &utils.UnknownError{}

}
