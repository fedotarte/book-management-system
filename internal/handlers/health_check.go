package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckHandler struct {
	log *logger.Logger
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{
		log: logger.GetLogger(),
	}
}

// HealthCheck      проверяет состояние сервиса
//
//	@Summary		проврка сервиса health-check
//	@Description	отдает 200, если все норм
//	@Tags			Utils
//	@Produce		json
//	@Success		200	{object}	dto.HealthCheckResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/health [get]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	h.log.Info("health check")
	res := dto.HealthCheckResponse{
		Status: "ok",
	}
	c.JSON(http.StatusOK, res)
}
