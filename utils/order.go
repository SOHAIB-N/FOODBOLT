package utils

import (
	"fmt"
	"time"
)

func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
} 