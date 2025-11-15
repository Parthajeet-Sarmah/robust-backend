package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"local/bomboclat-oidc-service/database"
	custom_errors "local/bomboclat-oidc-service/errors"
	custom_types "local/bomboclat-oidc-service/types"
	utils "local/bomboclat-oidc-service/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/golang-jwt/jwt/v5"
)

func (us *UserService) GetUserById(userId string, requiredFields string) (*custom_types.UserProfile, error) {

	var fields []string

	if strings.Contains(requiredFields, "profile") {
		fields = append(fields, "username")
	}
	if strings.Contains(requiredFields, "email") {
		fields = append(fields, "email")
	}

	fieldString := strings.Join(fields, ", ")

	data, err := database.FindUserByUUID(us.DBConn, context.Background(), userId, fieldString)

	if err != nil {
		return nil, err
	}

	profile := custom_types.UserProfile{
		Email:        data.Email,
		UserUUID:     data.UUID,
		Username:     data.Username,
		PasswordHash: data.PasswordHash,
	}

	return &profile, nil
}

func (us *UserService) Login(userDetails custom_types.UserLoginDetails) (*http.Cookie, error) {

	data, err := database.FindUserByEmailAndPasswordHash(us.DBConn, context.Background(),
		userDetails.Email, utils.HashToken256(userDetails.Password))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Invalid user details")
			return nil, custom_errors.UserNotFoundError(err)
		}
		return nil, err
	}

	if data == nil {
		log.Print("No user with this email!")
		return nil, custom_errors.UserNotFoundError(nil)
	}

	sessionID := uuid.New().String()

	userDetailsMap := map[string]string{
		"user_id":    data.UUID,
		"scope":      "deny",
		"expires_at": time.Now().Add(time.Minute * 5).Format(time.RFC3339),
	}

	if err := utils.SetValueToHash(us.RedisClient, "user_session:"+sessionID, userDetailsMap, time.Duration(time.Minute*5)); err != nil {
		log.Print(err)
		return nil, custom_errors.RedisSetHasError(err)
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

func (us *UserService) Logout(userCookie *http.Cookie) error {
	//remove redis session
	sessionId := userCookie.Value
	if err := utils.DeleteHashAll(us.RedisClient, "user_session:"+sessionId); err != nil {
		var appError *custom_errors.AppError
		if errors.As(err, &appError) && appError.Code == "REDIS_DELETE_ERROR" {
			return custom_errors.UserLogoutError(appError.Cause)
		}
	}

	fmt.Print("Success logout")

	return nil
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
		return custom_errors.UserAlreadyExistsError(nil)
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
		return nil, custom_errors.TokenParsingError(err)
	}

	if claims, ok := token.Claims.(*custom_types.CustomClaims); ok && token.Valid {

		userUUID := claims.Subject

		user, err := database.FindUserByUUID(us.DBConn, context.Background(), userUUID, "")

		userInfo := &custom_types.UserInfo{
			Username: user.Username,
			Email:    user.Email,
		}

		if err != nil {
			return nil, err
		}

		return userInfo, nil
	}

	return nil, custom_errors.Internal("Unknown error", nil)

}
