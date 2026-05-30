package handler

import (
	"errors"
	"net/http"
	"strconv"

	"GoTaskFlow/internal/model"
	"GoTaskFlow/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(ctx *gin.Context) {
	var request model.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), request)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) List(ctx *gin.Context) {
	tasks, err := h.service.List(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	task, err := h.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateStatus(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var request model.UpdateTaskStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.UpdateStatus(ctx.Request.Context(), id, request.Status)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *TaskHandler) Stats(ctx *gin.Context) {
	stats, err := h.service.Stats(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

func parseID(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		errorResponse(ctx, http.StatusBadRequest, "invalid task id")
		return 0, false
	}

	return id, true
}

func handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		errorResponse(ctx, http.StatusNotFound, "task not found")
	case errors.Is(err, service.ErrInvalidInput):
		errorResponse(ctx, http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrInvalidStatus):
		errorResponse(ctx, http.StatusBadRequest, "invalid task status")
	default:
		errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
}

func errorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message})
}
