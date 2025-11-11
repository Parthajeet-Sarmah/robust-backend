package custom_errors

import "net/http"

func MalformedJSONRequestError(status int, cause error) *AppError {
	return New("MALFORMED_JSON", "The JSON request body was malformed", status, cause)
}
func CouldNotConnectToDatabaseError(cause error) *AppError {
	return New("DB_CONN_FAILED", "Could not connect to database", http.StatusInternalServerError, cause)
}
func CouldNotFetchAuthCodeError(cause error) *AppError {
	return New("AUTH_CODE_FETCH_FAILED", "Auth code could not be fetched", http.StatusInternalServerError, cause)
}
func CouldNotFetchUserDataError(cause error) *AppError {
	return New("USER_DATA_FETCH_FAILED", "User data could not be fetched", http.StatusInternalServerError, cause)
}
func RedisCouldNotCreateClientError(cause error) *AppError {
	return New("REDIS_CLIENT_INIT_FAILED", "Could not initialise Redis client", http.StatusInternalServerError, cause)
}
func RedisGetHashError(cause error) *AppError {
	return New("REDIS_GET_ERROR", "The specified hash could not be fetched", http.StatusInternalServerError, cause)
}
func RedisSetHasError(cause error) *AppError {
	return New("REDIS_SET_ERROR", "The specified hash could not be set", http.StatusInternalServerError, cause)
}
func RedisGetHashNoResourceFoundError(cause error) *AppError {
	return New("REDIS_RESOURCE_NOT_FOUND", "The specified resource was not found", http.StatusNotFound, cause)
}
