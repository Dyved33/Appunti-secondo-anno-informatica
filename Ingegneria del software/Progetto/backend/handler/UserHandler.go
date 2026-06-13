package Handler

import (
	"context"
	"gialosai/backend/DTO"
	"gialosai/backend/domain"
	"gialosai/backend/service"
)

type UserHandler struct {
	Service *Service.UserService
}

func NewUserHandler(svc *Service.UserService) *UserHandler {
	return &UserHandler{Service: svc}
}

func (h *UserHandler) HandleLogin(ctx context.Context, req DTO.LoginRequestDTO) DTO.LoginResponseDTO {
	u, err := h.Service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return DTO.LoginResponseDTO{
			Success:      false,
			ErrorMessage: err.Error(),
		}
	}

	return DTO.LoginResponseDTO{
		Success:  true,
		Username: u.Username,
		Ruolo:    string(u.Ruolo),
	}
}

func (h *UserHandler) HandleGetSecurityQuestion(ctx context.Context, req DTO.RecoveryRequestDTO) (DTO.RecoveryResponseDTO, error) {
	q, err := h.Service.GetSecurityQuestion(ctx, req.Username)
	if err != nil {
		return DTO.RecoveryResponseDTO{}, err
	}
	return DTO.RecoveryResponseDTO{DomandaSicurezza: q}, nil
}

func (h *UserHandler) HandleResetPassword(ctx context.Context, req DTO.ResetPasswordRequestDTO) error {
	return h.Service.VerifySecurityAnswerAndResetPassword(ctx, req.Username, req.Risposta, req.NewPassword)
}