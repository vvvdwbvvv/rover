package model

import (
	"time"
)

// ContainerState 定義容器狀態存儲
type ContainerState struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
