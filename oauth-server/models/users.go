package models

import "time"

//	query := `CREATE TABLE IF NOT EXISTS users (
//		id SERIAL PRIMARY KEY,
//		username TEXT,
//		email TEXT,
//		password_hash TEXT,
//		uuid UUID DEFAULT gen_random_uuid(),
//		created_at TIMESTAMP DEFAULT now(),
//		updated_at TIMESTAMP DEFAULT now()
//	)`

type UserConsentsModel struct {
	UserId   string
	ClientId string
	Scopes   string
}

type UserDatabaseModel struct {
	Id           int `db:"-"`
	Username     string
	Email        string
	PasswordHash string `db:"-"`
	UUID         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
