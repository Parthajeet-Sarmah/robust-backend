package models

import "time"

type UserDatabaseModelInput struct {
	Username     string
	Email        string
	PasswordHash string
}

type UserDatabaseModel struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	UUID         string    `json:"uuid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
