package main

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config"
)

func main() {
	services := map[string]config.Service{
		"db":  {Name: "db", Image: "postgres:15", DependsOn: []string{}},
		"web": {Name: "web", Image: "nginx:latest", DependsOn: []string{"db"}},
		"app": {Name: "app", Image: "golang:latest", DependsOn: []string{"db", "web"}},
	}

	order, err := config.GetServiceStartupOrder(services)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Service startup order:", order)
}
