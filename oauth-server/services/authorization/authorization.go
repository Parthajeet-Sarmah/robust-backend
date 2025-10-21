package authorization

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/models"
	utils "local/bomboclat-oauth-server/utils"

	"github.com/jackc/pgx/v5"
)

func (as *AuthorizationService) AuthorizeUserAndGenerateCode(
	m models.AuthorizationRequestModelInput,
	userCookie *http.Cookie,
) (*string, error) {

	// NOTE: Check if client is registered with this service
	client, err := database.FindClientById(as.DBConn, context.Background(), m.ClientId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &utils.ClientNotFoundError{}
		}
		log.Print(err)
		return nil, err
	}

	if client == nil {
		return nil, &utils.ClientNotFoundError{}
	}

	if client.RedirectUri != nil && *client.RedirectUri != m.RedirectUri {
		log.Print("Redirect URI for the client does not match!")
		return nil, &utils.RedirectURIMismatchError{}
	}

	if userCookie == nil {
		return nil, &utils.UserNotLoggedInError{}
	}

	userKey := fmt.Sprintf("user_session:%s", userCookie.Value)
	res, err := utils.GetValueFromHash(as.RedisClient, userKey)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	doesUserSessionExist := res != nil
	scopeNotAllowed := res["scope"] != "deny"

	if !doesUserSessionExist {
		// Send back to controller with user not logged in error for redirection to /login
		return nil, &utils.UserNotLoggedInError{}
	}

	if !scopeNotAllowed {
		// Send back to controller with user scope denied error for redirection to /authorize/consent
		return nil, &utils.UserScopeDeniedError{}
	}

	randomBytes := make([]byte, 64)

	if _, err := rand.Read(randomBytes); err != nil {
		log.Print("Error while reading random bytes for generating code!")
		panic(err)
	}

	authCode := hex.EncodeToString(randomBytes)

	authCodeData := models.AuthCodeModelInput{
		Code:                authCode,
		UserId:              res["user_id"],
		ClientId:            *client.ClientId,
		Scopes:              m.Scope,
		RedirectUri:         *client.RedirectUri,
		Used:                false,
		CodeChallenge:       m.CodeChallenge,
		CodeChallengeMethod: m.CodeChallengeMethod,
	}

	if err := database.CreateAuthCodeEntry(as.DBConn, authCodeData); err != nil {
		return nil, err
	}

	url := *client.RedirectUri + "?code=" + authCode + "&state=" + m.State
	return &url, nil
}

func (as *AuthorizationService) AuthorizeConsent(m models.AuthorizationConsentModelInput, userCookie *http.Cookie) error {

	if m.Decision == "deny" {
		log.Print("Denied permission of data")
		return &utils.UserScopeDeniedError{}
	}

	if m.Decision == "allow" {

		//Get session id
		sessionId := userCookie.Value

		// TODO: Make a mechanism to check client_id and redirect_uri for double security
		res, err := utils.GetValueFromHash(as.RedisClient, "user_session:"+sessionId)

		if err != nil {
			log.Print(err)
			return &utils.RedisGetHashError{}
		}

		res["scope"] = "allow"

		utils.SetValueToHash(as.RedisClient, "user_session:"+sessionId, res)

		return nil
	}

	return errors.New("Wrong method")
}

func (as *AuthorizationService) GenerateToken() {

}

func (as *AuthorizationService) Introspect() {

}
