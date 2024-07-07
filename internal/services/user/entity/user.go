package entity

import (
	"github.com/gofrs/uuid"
)

type PosUser struct {
	Id            uuid.UUID `json:"id"`
	RestaurantId  uuid.UUID `json:"restaurant_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Position      string    `json:"position"`
	FirstPassword string    `json:"first_password"`
}

type BackofficeUser struct {
	Id            uuid.UUID `json:"id"`
	RestaurantId  uuid.UUID `json:"restaurant_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	FirstPassword string    `json:"first_password"`
}

type Customer struct {
	Id          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
}
