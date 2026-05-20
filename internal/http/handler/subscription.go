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

func (h *SubscriptionHandler) List(c *gin.Context) {
	filter, err := subscriptionFilterFromQuery(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	subscriptions, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	items := make([]subscriptionResponse, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		items = append(items, subscriptionResponseFromDomain(subscription))
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

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
	subscription.ID = id

	updated, err := h.service.Update(c.Request.Context(), subscription)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, subscriptionResponseFromDomain(updated))
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) TotalPrice(c *gin.Context) {
	filter, err := summaryFilterFromQuery(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	totalPrice, err := h.service.TotalPrice(c.Request.Context(), filter)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": totalPrice,
	})
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

func subscriptionFilterFromQuery(c *gin.Context) (repository.SubscriptionFilter, error) {
	var filter repository.SubscriptionFilter

	if value := c.Query("user_id"); value != "" {
		userID, err := uuid.Parse(value)
		if err != nil {
			return repository.SubscriptionFilter{}, err
		}

		filter.UserID = &userID
	}

	if value := c.Query("service_name"); value != "" {
		filter.ServiceName = &value
	}

	limit, err := intQuery(c, "limit")
	if err != nil {
		return repository.SubscriptionFilter{}, err
	}
	filter.Limit = limit

	offset, err := intQuery(c, "offset")
	if err != nil {
		return repository.SubscriptionFilter{}, err
	}
	filter.Offset = offset

	return filter, nil
}

func summaryFilterFromQuery(c *gin.Context) (repository.SubscriptionSummaryFilter, error) {
	from, err := domain.ParseYearMonth(c.Query("from"))
	if err != nil {
		return repository.SubscriptionSummaryFilter{}, err
	}

	to, err := domain.ParseYearMonth(c.Query("to"))
	if err != nil {
		return repository.SubscriptionSummaryFilter{}, err
	}

	filter := repository.SubscriptionSummaryFilter{
		From: from,
		To:   to,
	}

	if value := c.Query("user_id"); value != "" {
		userID, err := uuid.Parse(value)
		if err != nil {
			return repository.SubscriptionSummaryFilter{}, err
		}

		filter.UserID = &userID
	}

	if value := c.Query("service_name"); value != "" {
		filter.ServiceName = &value
	}

	return filter, nil
}

func intQuery(c *gin.Context, key string) (int, error) {
	value := c.Query(key)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
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
