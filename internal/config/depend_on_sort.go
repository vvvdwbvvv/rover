package config

import (
	"errors"
	"fmt"
)

func GetServiceStartupOrder(services map[string]Service) ([]string, error) {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	for name := range services {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	for name, service := range services {
		for _, dep := range service.DependsOn {
			if _, exists := services[dep]; !exists {
				return nil, fmt.Errorf("Service %s depends on unknown service %s", name, dep)
			}
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	order := []string{}
	queue := []string{}

	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(order) != len(services) {
		return nil, errors.New("Detected circular dependency in depends_on")
	}

	return order, nil
}
