package handler

import (
	"agnos-gin/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service *service.PatientService
}

func NewPatientHandler(s *service.PatientService) *PatientHandler {
	return &PatientHandler{service: s}
}

func (h *PatientHandler) Search(c *gin.Context) {
	hospital, exists := c.Get("hospital_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filters := make(map[string]interface{})
	if v := c.Query("national_id"); v != "" { filters["national_id"] = v }
	if v := c.Query("passport_id"); v != "" { filters["passport_id"] = v }
	if v := c.Query("first_name"); v != "" { filters["first_name"] = v }

	patients, err := h.service.Search(filters, hospital.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": patients})
}