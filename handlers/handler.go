package handlers

import (
	"github.com/mahi7000/bookstore-go-app/internal/hasura"
)

type ApiConfig struct {
	DB *hasura.Client
}

func NewHandlerConfig(hasuraClient *hasura.Client) *ApiConfig {
	return &ApiConfig{
		DB: hasuraClient,
	}
}