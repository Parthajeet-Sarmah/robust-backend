package models

import "time"

type ConsentModel struct {
	Id        int
	UserId    string
	ClientId  string
	Scopes    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
