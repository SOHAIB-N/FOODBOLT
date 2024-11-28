package models

import (
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

type Order struct {
	gorm.Model
	UserID      uint
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	TotalAmount float64
	Status      string // pending, confirmed, preparing, ready, delivered
	PaymentID   string
	IsDelivery  bool
	RoomNumber  string
	OrderNumber string `gorm:"unique"`
}

type OrderItem struct {
	gorm.Model
	OrderID    uint
	MenuItemID uint
	MenuItem   MenuItem
	Quantity   int
	Price      float64
	TotalPrice float64
}

type Payment struct {
	gorm.Model
	OrderID       uint
	Amount        float64
	TransactionID string
	Status        string
	PaymentMethod string
}