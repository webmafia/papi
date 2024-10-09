package main

import "time"

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" flags:"required"`
	TimeCreated time.Time `json:"timeCreated"`
	Awesome     bool      `json:"awesome" title:"Whether the user is awesome" flags:"readonly,nullable"`
}
