package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Graph[K comparable, N any, E any] struct {
	nodes map[K]N
	edges map[K]map[K]E
}

func NewGraph[K comparable, N any, E any]() *Graph[K, N, E] {
	nodes := make(map[K]N)
	edges := make(map[K]map[K]E)
	return &Graph[K, N, E]{nodes, edges}
}

func parse(file *os.File) (*Graph[string, uint64, uint64], error) {
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile("\\AValve ([A-Z]+).+rate=(\\d+);.+ to valves? ([A-Z]+(?:,\\s*[A-Z]+)*)\\z")

	graph := NewGraph[string, uint64, uint64]()

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindAllStringSubmatch(line, -1)[0]

		id := match[1]
		rate, err := strconv.ParseUint(match[2], 10, 64)
		neighbors := strings.Split(match[3], ", ")

		if err != nil {
			return nil, err
		}

		graph.nodes[id] = rate

		edges := make(map[string]uint64)

		for _, nid := range neighbors {
			edges[nid] = 1
		}

		graph.edges[id] = edges
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return graph, nil
}

func main() {
	start := "AA"

	graph, err := parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input", err)
	}

	fmt.Printf("%v\n%v\n", graph.nodes, graph.edges)
}
