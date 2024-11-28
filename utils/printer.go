package utils

import (
	"fmt"
	"food-court/models"
	"time"
)

func PrintReceipt(order models.Order) error {
	// This is a placeholder for actual thermal printer implementation
	// You would need to implement the specific printer driver integration here
	receipt := fmt.Sprintf(`
	=========================
	FOOD COURT ORDER RECEIPT
	=========================
	Order #: %s
	Date: %s
	--------------------------
	`, order.OrderNumber, time.Now().Format("2006-01-02 15:04:05"))

	for _, item := range order.Items {
		itemLine := fmt.Sprintf("%dx %s - $%.2f\n", item.Quantity, item.MenuItem.Name, item.TotalPrice)
		receipt += itemLine
	}

	receipt += fmt.Sprintf(`
	--------------------------
	Total: $%.2f
	=========================
	`, order.TotalAmount)

	// Here you would send the receipt to the actual printer
	fmt.Println(receipt) // For demonstration purposes
	return nil
}