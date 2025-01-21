package main

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Parsed Config: %+v\n", cfg)
}
