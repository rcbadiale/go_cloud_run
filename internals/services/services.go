package services

import "github.com/rcbadiale/go-cloud-run/internals"

type BaseHttpService struct {
	Client internals.HTTPClient
}
