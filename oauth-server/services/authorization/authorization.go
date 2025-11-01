package authorization

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"
	utils "local/bomboclat-oauth-server/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func (as *AuthorizationService) AuthorizeUserAndGenerateCode(
	m custom_types.AuthorizationRequestModelInput,
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

	authCodeData := custom_types.AuthCodeModelInput{
		Code:                authCode,
		UserId:              res["user_id"],
		ClientId:            *client.ClientId,
		Scopes:              m.Scope,
		RedirectUri:         *client.RedirectUri,
		Used:                false,
		CodeChallenge:       m.CodeChallenge,
		CodeChallengeMethod: m.CodeChallengeMethod,
	}

	if err := database.CreateAuthCodeEntry(as.DBConn, &authCodeData); err != nil {
		return nil, err
	}

	url := *client.RedirectUri + "?code=" + authCode + "&state=" + m.State
	return &url, nil
}

func (as *AuthorizationService) AuthorizeConsent(m custom_types.AuthorizationConsentModelInput, userCookie *http.Cookie) error {

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

func (as *AuthorizationService) GenerateToken(m *custom_types.TokenModelInput) (*custom_types.TokenResponse, error) {

	switch m.GrantType {
	case "authorization_code":
		//Validate the client with the respective ClientId
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

		//Validate the code with the ClientId and Redirect_Uri
		codeData, err := database.GetAuthCode(as.DBConn, m.Code)

		if err != nil {
			log.Print(err)
			return nil, &utils.CouldNotFetchAuthCode{}
		}

		if codeData.ClientId != m.ClientId {
			return nil, &utils.ClientIdMismatchError{}
		} else if codeData.RedirectUri != m.RedirectUri {
			return nil, &utils.RedirectURIMismatchError{}
		} else if time.Now().UTC().Compare(codeData.ExpiresAt) == 1 {
			return nil, &utils.ExpiredAuthCodeError{}
		}

		//Verify the PKCE challenge
		var codeChallenge string

		if m.CodeChallengeMethod == "S256" {
			hash := sha256.Sum256([]byte(m.CodeVerifier))
			codeChallenge = base64.RawURLEncoding.EncodeToString(hash[:])
		} else {
			codeChallenge = m.CodeVerifier
		}

		if codeChallenge != codeData.CodeChallenge {
			return nil, &utils.CodeChallengeDoesNotMatchError{}
		}

		// NOTE: Generate access and refresh token (optionally) and ID token (if OIDC)
		expiresAt := time.Now().UTC().Add(10 * time.Minute) // 10 minutes
		tokenClaims := jwt.MapClaims{
			"sub":   codeData.UserId,
			"aud":   m.ClientId,
			"exp":   expiresAt.Unix(),
			"iat":   time.Now().Unix(),
			"scope": codeData.Scopes,
			"iss":   os.Getenv("OIDC_BASE_URL"),
		}

		// NOTE: Get RSA private key
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_RSA_PRIVATE_KEY")))
		if err != nil {
			return nil, err
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
		accessToken, err := jwtToken.SignedString(key)
		if err != nil {
			return nil, err
		}

		// NOTE: Save access token
		if err := database.InsertAccessToken(as.DBConn, &models.AccessTokenModel{
			TokenHash: utils.HashToken256(accessToken),
			ClientId:  m.ClientId,
			UserId:    codeData.UserId,
			Scopes:    codeData.Scopes,
			ExpiresAt: expiresAt,
		}); err != nil {
			return nil, err
		}

		randomBytes := make([]byte, 64)

		if _, err := rand.Read(randomBytes); err != nil {
			log.Print("Error while reading random bytes for generating code!")
			panic(err)
		}

		refreshToken := hex.EncodeToString(randomBytes)

		// NOTE: Save refresh token
		if err := database.InsertRefreshToken(as.DBConn, &models.RefreshTokenModel{
			TokenHash: utils.HashToken256(refreshToken),
			ClientId:  m.ClientId,
			UserId:    codeData.UserId,
			Scopes:    codeData.Scopes,
			ExpiresAt: time.Now().UTC().Add(30 * 24 * time.Hour),
		}); err != nil {
			return nil, err
		}

		if err := database.UpdateAuthCodeEntryUsedStatus(as.DBConn, m.Code); err != nil {
			return nil, &utils.AuthCodeUsedUpdateError{}
		}

		// NOTE: Mark auth code as used
		if err := database.UpdateAuthCodeEntryUsedStatus(as.DBConn, m.Code); err != nil {
			return nil, &utils.AuthCodeUsedUpdateError{}
		}

		return &custom_types.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(time.Until(expiresAt).Seconds()),
		}, nil

	case "refresh_token":
		tokenData, err := database.FindRefreshToken(as.DBConn, utils.HashToken256(m.RefreshToken))

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, &utils.RefreshTokenNotFoundError{
					Status: http.StatusNotFound,
					Msg:    "Refresh token not found",
				}
			}
			return nil, err
		}

		if tokenData == nil {
			return nil, &utils.RefreshTokenNotFoundError{
				Status: http.StatusNotFound,
				Msg:    "Refresh token not found",
			}
		}

		if tokenData.ClientId != m.ClientId {
			return nil, &utils.ClientIdMismatchError{}
		}

		if time.Now().UTC().Compare(tokenData.ExpiresAt) == 1 {
			return nil, &utils.ExpiredRefreshTokenError{}
		}

		expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes
		tokenClaims := jwt.MapClaims{
			"sub":   utils.HashToken256(tokenData.UserId),
			"aud":   tokenData.ClientId,
			"exp":   expiresAt.Unix(),
			"iat":   time.Now().Unix(),
			"scope": tokenData.Scopes,
		}

		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_RSA_PRIVATE_KEY")))
		if err != nil {
			return nil, err
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
		accessToken, err := jwtToken.SignedString(key)
		if err != nil {
			return nil, err
		}

		if err := database.UpdateAccessToken(as.DBConn, &models.AccessTokenModel{
			TokenHash: utils.HashToken256(accessToken),
			ClientId:  tokenData.ClientId,
			UserId:    tokenData.UserId,
			Scopes:    tokenData.Scopes,
			ExpiresAt: expiresAt,
		}); err != nil {
			return nil, err
		}

		// NOTE: Rotate refresh token and update it in DB
		randomBytes := make([]byte, 64)

		if _, err := rand.Read(randomBytes); err != nil {
			log.Print("Error while reading random bytes for generating code!")
			panic(err)
		}

		refreshToken := hex.EncodeToString(randomBytes)

		err = database.UpdateRefreshTokenEntry(as.DBConn, utils.HashToken256(m.RefreshToken), utils.HashToken256(refreshToken))
		if err != nil {
			return nil, err
		}

		return &custom_types.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(time.Until(expiresAt).Seconds()),
		}, nil

	default:
		return nil, &utils.InvalidGrantType{}
	}
}

// Idempotent call
func (as *AuthorizationService) RevokeToken(m *models.RevokeTokenModel) error {

	switch m.TokenTypeHint {
	case "access_token":
		tokenData, err := database.FindAccessToken(as.DBConn, utils.HashToken256(m.Token))

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}

		if tokenData == nil {
			return nil
		}

		if !tokenData.Revoked {
			database.RevokeAccessToken(as.DBConn, tokenData.TokenHash)
		}
	}

	return nil
}
