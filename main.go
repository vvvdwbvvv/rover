/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config" // ✅ 确保正确导入
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Parsed Config: %+v\n", cfg)
}
