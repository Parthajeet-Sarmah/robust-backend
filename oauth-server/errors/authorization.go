package custom_errors

import "net/http"

func UnauthorizedClientError(cause error) *AppError {
	return New("UNAUTHORIZED_CLIENT", "The client is unauthorized", http.StatusUnauthorized, cause)
}
func ExpiredAuthCodeError(cause error) *AppError {
	return New("AUTH_CODE_EXPIRED", "The auth code has expired", http.StatusGone, cause)
}
func ExpiredRefreshTokenError(cause error) *AppError {
	return New("REFRESH_TOKEN_EXPIRED", "The refresh token has expired", http.StatusGone, cause)
}
func AuthCodeUsedUpdateError(cause error) *AppError {
	return New("AUTH_CODE_UPDATE_ERROR", "Error in updating auth code entry", http.StatusInternalServerError, cause)
}
func CodeChallengeDoesNotMatchError(cause error) *AppError {
	return New("CHALLENGE_MISMATCH", "Code challenge does not match", http.StatusForbidden, cause)
}
func InvalidGrantTypeError(cause error) *AppError {
	return New("INVALID_GRANT", "Invalid grant_type provided", http.StatusBadRequest, cause)
}
func TokenParsingError(cause error) *AppError {
	return New("TOKEN_PARSING_ERROR", "Token could not be parsed", http.StatusInternalServerError, cause)
}
func NoAccessTokenFoundError(cause error) *AppError {
	return New("ACCESS_TOKEN_NOT_FOUND", "No access token found", http.StatusNotFound, cause)
}
func RefreshTokenNotFoundError(cause error) *AppError {
	return New("REFRESH_TOKEN_NOT_FOUND", "No refresh token found", http.StatusNotFound, cause)
}
