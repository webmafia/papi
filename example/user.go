package main

import "time"

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TimeCreated time.Time `json:"timeCreated"`
}
