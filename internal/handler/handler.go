package handler

import "pr_task/internal/service"

type Handler struct {
	Service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{
		Service: service,
	}
}
