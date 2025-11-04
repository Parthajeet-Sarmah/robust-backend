package custom_types

import "github.com/jackc/pgx/v5/pgxpool"

type UserDetails struct {
	Email        string
	PasswordHash string
}

type Postgres struct {
	DB *pgxpool.Pool
}

type UserDatabaseModelInput struct {
	Username     string
	Email        string
	PasswordHash string
}

type UserProfile struct {
	UserUUID        string `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	PasswordHash    string `json:"password_hash,omitempty"`
	IsEmailVerified bool   `json:"email_verified,omitempty"`
	ProfilePic      string `json:"profile_pic,omitempty"`
}
