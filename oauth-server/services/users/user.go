package users

import (
	"context"
	"errors"
	"log"
	"net/http"

	"local/bomboclat-oauth-server/database"
	custom_types "local/bomboclat-oauth-server/types"
	utils "local/bomboclat-oauth-server/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (us *UserService) Login(userDetails custom_types.UserDetails) (*http.Cookie, error) {

	data, err := database.FindUserByEmailAndPasswordHash(us.DBConn, context.Background(),
		userDetails.Email, userDetails.PasswordHash)

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
