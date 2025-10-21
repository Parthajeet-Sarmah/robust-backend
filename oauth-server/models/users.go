package models

import "time"

type UserDatabaseModelInput struct {
	Username     string
	Email        string
	PasswordHash string
}

//	query := `CREATE TABLE IF NOT EXISTS users (
//		id SERIAL PRIMARY KEY,
//		username TEXT,
//		email TEXT,
//		password_hash TEXT,
//		uuid UUID DEFAULT gen_random_uuid(),
//		created_at TIMESTAMP DEFAULT now(),
//		updated_at TIMESTAMP DEFAULT now()
//	)`

type UserDatabaseModel struct {
	Id           int
	Username     string
	Email        string
	PasswordHash string
	UUID         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
