package service

import (
	"github.com/google/uuid"
	"time"
)

type CreateMinerInput struct {
	Name   string `json:"name"`
	Energy int    `json:"energy"`
	Age    int    `json:"age"`
}

type CreateMinerOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Energy    int       `json:"energy"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}
