package graphql

import "github.com/m-mizutani/alertchain/pkg/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	svc *service.Services
}

func NewResolver(svc *service.Services) *Resolver {
	return &Resolver{
		svc: svc,
	}
}
