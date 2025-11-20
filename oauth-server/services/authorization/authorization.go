package authorization

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"local/bomboclat-oauth-server/database"
	custom_errors "local/bomboclat-oauth-server/errors"
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

	if m.State == "" {
		return nil, custom_errors.Internal("No state was provided", nil)
	}

	if userCookie != nil {
		res, err := utils.GetValueFromHash(as.RedisClient, "user_session:"+userCookie.Value)

		if err != nil {
			return nil, custom_errors.RedisGetHashError(err)
		}

		expiryStr, exists := res["expires_at"]
		if !exists {
			return nil, custom_errors.Internal("Session has no expiry", nil)
		}

		expiry, err := time.Parse(time.RFC3339, expiryStr)
		if err != nil || time.Now().After(expiry) {
			return nil, custom_errors.Internal("Session has expired", nil)
		}
	}

	// NOTE: Check if client is registered with this service
	client, err := database.FindClientById(as.DBConn, context.Background(), m.ClientId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_errors.ClientNotFoundError(err)
		}
		return nil, err
	}

	if client.RedirectUri == "" {
		return nil, custom_errors.ClientNotFoundError(nil)
	}

	// TODO: Multiple URL check, open redirects, HTTPS check
	if client.RedirectUri != "" && client.RedirectUri != m.RedirectUri {
		return nil, custom_errors.RedirectURIMismatchError(nil)
	}

	if userCookie == nil {
		return nil, custom_errors.UserNotLoggedInError(nil)
	}

	userKey := fmt.Sprintf("user_session:%s", userCookie.Value)
	res, err := utils.GetValueFromHash(as.RedisClient, userKey)

	if err != nil {
		return nil, err
	}

	doesUserSessionExist := res != nil

	// NOTE: Check for scopes in the 'consents' table
	consent, err := database.FindUserConsent(as.DBConn, context.Background(), res["user_id"], m.ClientId)
	isScopePresentAndEqual := consent != nil && consent.Scopes == m.Scope

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	scopeAllowed := res["scope"] != "deny" && isScopePresentAndEqual

	if !doesUserSessionExist {
		// Send back to controller with user not logged in error for redirection to /login
		return nil, custom_errors.UserNotLoggedInError(nil)
	}

	if !scopeAllowed {
		// Send back to controller with user scope denied error for redirection to /authorize/consent
		return nil, custom_errors.UserScopeDeniedError(nil)
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
		ClientId:            client.ClientId,
		Scopes:              m.Scope,
		RedirectUri:         client.RedirectUri,
		Used:                false,
		CodeChallenge:       m.CodeChallenge,
		CodeChallengeMethod: m.CodeChallengeMethod,
	}

	if err := database.CreateAuthCodeEntry(as.DBConn, &authCodeData); err != nil {
		return nil, err
	}

	url := client.RedirectUri + "?code=" + authCode + "&state=" + m.State
	return &url, nil
}

func (as *AuthorizationService) AuthorizeConsent(m custom_types.AuthorizationConsentModelInput, userCookie *http.Cookie) error {

	if m.Decision == "deny" {
		return custom_errors.UserScopeDeniedError(nil)
	}

	if m.Decision == "allow" {

		//Get session id
		sessionId := userCookie.Value

		// TODO: Make a mechanism to check client_id and redirect_uri for double security
		res, err := utils.GetValueFromHash(as.RedisClient, "user_session:"+sessionId)

		if err != nil {
			var appErr *custom_errors.AppError
			if errors.As(err, &appErr) && appErr.Code == "REDIS_RESOURCE_NOT_FOUND" {
				return custom_errors.New("USER_SESSION_NOT_FOUND", "This user session was not found", http.StatusUnauthorized, err)
			} else {
				return err
			}
		}

		res["scope"] = "allow"

		utils.SetValueToHash(as.RedisClient, "user_session:"+sessionId, res)

		consent := custom_types.UserConsentInput{
			UserId:   res["user_id"],
			ClientId: m.ClientId,
			Scopes:   m.Scope,
		}

		// NOTE: Persist consent in 'consents' table for future use
		database.InsertUserConsent(as.DBConn, context.Background(), &consent)

		return nil
	}

	return errors.New("Wrong method")
}

func (as *AuthorizationService) GenerateToken(m *custom_types.TokenModelInput, authMethod string) (*custom_types.TokenResponse, error) {

	switch m.GrantType {
	case "authorization_code":
		//Validate the client with the respective ClientId
		client, err := database.FindClientById(as.DBConn, context.Background(), m.ClientId)

		if authMethod != client.TokenEndpointAuthMethod {
			return nil, custom_errors.TokenEndpointAuthMethodMismatch(nil)
		}

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, custom_errors.ClientNotFoundError(err)
			}
			// TODO: Better errors here???
			return nil, custom_errors.ClientNotFoundError(err)
		}

		if client.RedirectUri == "" {
			return nil, custom_errors.ClientNotFoundError(nil)
		}

		if client.RedirectUri != "" && client.RedirectUri != m.RedirectUri {
			log.Print("Redirect URI for the client does not match!")
			return nil, custom_errors.RedirectURIMismatchError(nil)
		}

		//Validate the code with the ClientId and Redirect_Uri
		codeData, err := database.GetAuthCode(as.DBConn, m.Code)

		if err != nil {
			log.Print(err)
			return nil, err
		}

		if codeData.ClientId != m.ClientId {
			return nil, custom_errors.ClientIdMismatchError(nil)
		} else if codeData.RedirectUri != m.RedirectUri {
			return nil, custom_errors.RedirectURIMismatchError(nil)
		} else if time.Now().UTC().Compare(codeData.ExpiresAt) == 1 {
			return nil, custom_errors.ExpiredAuthCodeError(nil)
		} else if codeData.Used {
			return nil, custom_errors.AuthCodeAlreadyUsedError(nil)
		}

		//Verify the PKCE challenge
		var codeChallenge string

		if m.CodeChallengeMethod == "S256" {
			codeChallenge = utils.HashToken256(m.CodeVerifier)
		}

		if codeChallenge != codeData.CodeChallenge {
			return nil, custom_errors.CodeChallengeDoesNotMatchError(nil)
		}

		// NOTE: Do checks necessary for provided token_endpoint_auth_method
		if client.TokenEndpointAuthMethod == "private_key_jwt" {
			// NOTE: Break down JWT, get public key from client
			key, err := utils.BuildRSAPublicKey([]byte(*client.Jwks), "test-key-001")

			if err != nil {
				return nil, err
			}

			token, err := jwt.ParseWithClaims(m.ClientAssertion, &custom_types.CustomClaims{}, func(token *jwt.Token) (any, error) {
				return key, nil
			})

			// NOTE: Validate JWT claims
			if claims, ok := token.Claims.(*custom_types.CustomClaims); ok && token.Valid {
				iss := claims.Issuer
				sub := claims.Subject
				//aud := claims.Audience
				jti := claims.ID
				exp := claims.ExpiresAt

				if iss != client.ClientId || sub != client.ClientId || jti == "0" || exp == 0 {
					return nil, custom_errors.TokenParsingError(nil)
				}
			}
		}

		var idToken string
		// NOTE: Checking for OpenID scope and generate ID token if exists
		if strings.Contains(codeData.Scopes, "openid") {

			var user custom_types.UserProfile
			// Call OIDC service to fetch user data
			oidcBaseUrl := os.Getenv("OIDC_BASE_URL")
			client := &http.Client{}

			// TODO: Add a internal service token
			fields := strings.Split(codeData.Scopes, " ")
			var fieldsFiltered []string

			for _, val := range fields {
				if !strings.Contains(val, "openid") {
					fieldsFiltered = append(fieldsFiltered, val)
				}
			}

			fieldString := "?fields=" + strings.Join(fieldsFiltered, ",")

			req, err := http.NewRequest("GET", oidcBaseUrl+"/users/id/"+codeData.UserId+fieldString, nil)
			//req.Header.Set("Authorization", "Bearer internal")

			if err != nil {
				return nil, err
			}

			resp, err := client.Do(req)

			if err != nil {
				return nil, err
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(body, &user)

			if err != nil {
				return nil, err
			}

			expiresAt := time.Now().UTC().Add(10 * time.Minute) // 10 minutes

			// NOTE: Create ID token if user is not null
			var pidTokenClaims jwt.MapClaims

			pidTokenClaims = jwt.MapClaims{
				"iss":   os.Getenv("OIDC_BASE_URL"),
				"sub":   user.UserUUID,
				"aud":   m.ClientId,
				"exp":   expiresAt.Unix(),
				"iat":   time.Now().Unix(),
				"email": user.Email,
				"name":  user.Username,
			}

			key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_RSA_PRIVATE_KEY")))
			if err != nil {
				return nil, err
			}

			jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, pidTokenClaims)
			tokenString, err := jwtToken.SignedString(key)
			if err != nil {
				return nil, err
			}

			idToken = tokenString

		}

		expiresAt := time.Now().UTC().Add(10 * time.Minute) // 10 minutes

		// NOTE: Generate access and refresh token (optionally) and ID token (if OIDC)
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

		// TODO: Conditional wrap logic, wrap the error first in DB
		// NOTE: Mark auth code as used
		if err := database.UpdateAuthCodeEntryUsedStatus(as.DBConn, m.Code); err != nil {
			return nil, custom_errors.AuthCodeUsedUpdateError(err)
		}

		return &custom_types.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			IdToken:      idToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(time.Until(expiresAt).Seconds()),
		}, nil

	case "refresh_token":
		tokenData, err := database.FindRefreshToken(as.DBConn, utils.HashToken256(m.RefreshToken))

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, custom_errors.RefreshTokenNotFoundError(err)
			}
			return nil, custom_errors.RefreshTokenNotFoundError(err)
		}

		if tokenData == nil {
			return nil, custom_errors.RefreshTokenNotFoundError(nil)
		}

		if tokenData.ClientId != m.ClientId {
			return nil, custom_errors.RefreshTokenNotFoundError(nil)
		}

		if time.Now().UTC().Compare(tokenData.ExpiresAt) == 1 {
			return nil, custom_errors.ExpiredRefreshTokenError(nil)
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
		return nil, custom_errors.InvalidGrantTypeError(nil)
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
