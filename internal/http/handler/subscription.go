package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gokost710/subscription-service/internal/domain"
	"github.com/gokost710/subscription-service/internal/repository"
	"github.com/gokost710/subscription-service/internal/service"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service service.SubscriptionProvider
}

func NewSubscriptionHandler(service service.SubscriptionProvider) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var request subscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	subscription, err := request.toDomain()
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.service.Create(c.Request.Context(), subscription)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, subscriptionResponseFromDomain(created))
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	subscription, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, subscriptionResponseFromDomain(subscription))
}

func (h *SubscriptionHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidSubscription):
		writeError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrNotFound):
		writeError(c, http.StatusNotFound, "subscription not found")
	default:
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}

type subscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func (r subscriptionRequest) toDomain() (domain.Subscription, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return domain.Subscription{}, err
	}

	startDate, err := domain.ParseYearMonth(r.StartDate)
	if err != nil {
		return domain.Subscription{}, err
	}

	var endDate *domain.YearMonth
	if r.EndDate != nil {
		parsedEndDate, err := domain.ParseYearMonth(*r.EndDate)
		if err != nil {
			return domain.Subscription{}, err
		}

		endDate = &parsedEndDate
	}

	return domain.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

type subscriptionResponse struct {
	ID          int64   `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func subscriptionResponseFromDomain(subscription domain.Subscription) subscriptionResponse {
	var endDate *string
	if subscription.EndDate != nil {
		value := subscription.EndDate.String()
		endDate = &value
	}

	return subscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID.String(),
		StartDate:   subscription.StartDate.String(),
		EndDate:     endDate,
	}
}

func parseIDParam(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}
