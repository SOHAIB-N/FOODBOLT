package handlers

import (
	"encoding/json"
	"food-court/models"
	"food-court/utils"
	"net/http"
	

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Generate unique order number
	order.OrderNumber = utils.GenerateOrderNumber()
	order.Status = "pending"

	// Calculate total amount
	var totalAmount float64
	for i, item := range order.Items {
		var menuItem models.MenuItem
		if err := h.DB.First(&menuItem, item.MenuItemID).Error; err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item")
			return
		}
		order.Items[i].Price = menuItem.Price
		order.Items[i].TotalPrice = menuItem.Price * float64(item.Quantity)
		totalAmount += order.Items[i].TotalPrice
	}
	order.TotalAmount = totalAmount

	if err := h.DB.Create(&order).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating order")
		return
	}

	// Notify kitchen staff via WebSocket
	utils.BroadcastNewOrder(order)

	utils.RespondWithJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var order models.Order
	if err := h.DB.Preload("Items.MenuItem").First(&order, orderID).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", statusUpdate.Status).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating order status")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order status updated"})
}