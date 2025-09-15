package handler

import "github.com/novaru/billing-service/internal/app/service"

type Handlers struct {
	User *UserHandler
	Plan *PlanHandler
}

func New(
	userService service.UserService,
	planService service.PlanService,
) *Handlers {
	return &Handlers{
		User: NewUserHandler(userService),
		Plan: NewPlanHandler(planService),
	}
}
