package dimodel

import (
	"polling_websocket/pkg/domain/repos"
	"polling_websocket/pkg/interfaces/controllers"
)

type Dependencies struct {
	AuthService       *repos.AuthService
	AuthController    *controllers.AuthController
	PollingController *controllers.PollingController
}
