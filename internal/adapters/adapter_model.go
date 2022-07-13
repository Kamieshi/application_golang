package adapters

import "app/internal/service"

type ServicesApp struct {
	AuthService   *service.AuthService
	EntityService *service.EntityService
	ImageService  *service.ImageService
	UserService   *service.UserService
}
