package models

import "time"

type AuthorizationRequestModelInput struct {
	ResponseType        string
	ClientId            string
	RedirectUri         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type AuthorizationConsentModelInput struct {
	ClientId    string
	Scope       string
	Decision    string
	RedirectUri string
}

type AuthCodeModelInput struct {
	Id                  string
	Code                string
	UserId              string
	ClientId            string
	RedirectUri         string
	Scopes              string
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool
}

type AuthCodeModel struct {
	Id                  string
	Code                string
	UserId              string
	ClientId            string
	RedirectUri         string
	Scopes              string
	ExpiresAt           time.Time
	CreatedAt           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool
}
