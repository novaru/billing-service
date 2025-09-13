package handler

import "github.com/novaru/billing-service/internal/app/service"

type Handlers struct {
	User *UserHandler
}

func New(
	userService service.UserService,
) *Handlers {
	return &Handlers{
		User: NewUserHandler(userService),
	}
}
