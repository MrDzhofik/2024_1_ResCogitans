package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type RegistrationHandler struct {
	sessionUseCase usecase.SessionInterface
	userUseCase    usecase.UserUseCaseInterface
}

func NewRegistrationHandler(sessionUseCase usecase.SessionInterface, userUseCase usecase.UserUseCaseInterface) *RegistrationHandler {
	return &RegistrationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	sessionID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		httpError := httperrors.UnwrapHttpError(err)
		if httpError.Message != "Session not found" && httpError.Message != "Cookie not found" {
			return entities.UserResponse{}, httpError
		}
	}

	if sessionID != 0 {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "User is already authorized")
	}

	taken, err := h.userUseCase.IsEmailTaken(username)
	if err != nil {
		return entities.UserResponse{}, err
	}
	if taken {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "Username is taken")
	}

	if err := h.userUseCase.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, err
	}

	err = h.userUseCase.CreateUser(username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	user, err := h.userUseCase.GetUserByEmail(username)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.sessionUseCase.CreateSession(responseWriter, user.ID)
	if err != nil {
		return entities.UserResponse{}, err
	}

	userResponse := entities.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}
