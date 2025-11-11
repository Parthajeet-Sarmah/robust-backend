package custom_errors

import "net/http"

func UserAlreadyExistsError(cause error) *AppError {
	return New("USER_ALREADY_EXISTS", "The user already exists", http.StatusConflict, cause)
}
func UserNotFoundError(cause error) *AppError {
	return New("USER_NOT_FOUND", "The user was not found", http.StatusNotFound, cause)
}
func UserNotLoggedInError(cause error) *AppError {
	return New("USER_NOT_LOGGED_IN", "The user was not logged in", http.StatusUnauthorized, cause)
}
func UserScopeDeniedError(cause error) *AppError {
	return New("USER_SCOPE_DENIED", "The user denied permission for this scope", http.StatusForbidden, cause)
}
