package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"local/bomboclat-oauth-server/database"
	custom_errors "local/bomboclat-oauth-server/errors"
	"local/bomboclat-oauth-server/utils"
	"net/http"
	"strings"
	"time"
)

func (ms *MiddlewareServiceContainer) AuthorizeClient(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return "", custom_errors.UnauthorizedClientError(nil)
	}

	encoded := strings.TrimPrefix(auth, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)

	// TODO: Make an error for this later
	if err != nil {
		return "", custom_errors.Internal("An unexpected error occured", err)
	}

	parts := strings.SplitN(string(decodedBytes), ":", 2)
	if len(parts) != 2 {
		return "", custom_errors.Internal("Malformed authorization header", nil)
	}

	clientId, clientSecretHash := parts[0], utils.HashToken256(parts[1])
	// Check for clientId and secret in Redis first then to the database
	_, err = utils.GetValueFromHash(ms.RedisClient, fmt.Sprintf("%s:%s", clientId, clientSecretHash))

	if err != nil {
		//Check database
		client, err := database.FindClientByIdAndSecretHash(ms.DBConn, context.Background(), clientId, clientSecretHash)
		if err != nil {
			return "", custom_errors.UnauthorizedClientError(nil)
		}

		id := *client.ClientId
		secret := *client.ClientSecretHash

		data := map[string]string{
			"client_id":          id,
			"client_secret_hash": secret,
		}

		utils.SetValueToHash(ms.RedisClient, fmt.Sprintf("%s:%s", id, secret), data, 72*time.Hour)
	}

	return clientId, nil
}
