package auth

import "context"

type authService interface {
	AuthService_Login(ctx context.Context, username, password string) (string, error)
}

type AuthController struct {
	Service authService
}

func NewAuthController(svc authService) *AuthController {
	return &AuthController{Service: svc}
}

func (c *AuthController) AuthController_Login(ctx context.Context, req LoginRequestDTO) (LoginResponseDTO, error) {
	token, err := c.Service.AuthService_Login(ctx, req.Username, req.Password)
	if err != nil {
		return LoginResponseDTO{}, err
	}
	return LoginResponseDTO{Token: token}, nil
}