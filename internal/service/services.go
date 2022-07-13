package service

type ServicesApp struct {
	AuthService   *AuthService
	EntityService *EntityService
	ImageService  *ImageService
	UserService   *UserService
}
