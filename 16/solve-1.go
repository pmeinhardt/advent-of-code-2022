package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Graph struct {
	index map[string]int
	nodes []int
	edges [][]int
	none  int
}

func NewGraph() *Graph {
	index := make(map[string]int, 0)
	nodes := make([]int, 0)
	edges := make([][]int, 0)
	return &Graph{index, nodes, edges, math.MaxInt}
}

func (g *Graph) Size() int {
	return len(g.index)
}

func (g *Graph) Keys() []string {
	keys := make([]string, len(g.index))

	for k, i := range g.index {
		keys[i] = k
	}

	return keys
}

func (g *Graph) HasNode(key string) bool {
	_, present := g.index[key]
	return present
}

func (g *Graph) AddNode(key string, n int) error {
	_, present := g.index[key]

	if present {
		return fmt.Errorf("Node with key %v already exists", key)
	}

	size := len(g.index)

	g.index[key] = size
	g.nodes = append(g.nodes, n)

	for i, e := range g.edges {
		g.edges[i] = append(e, g.none)
	}

	edges := make([]int, size+1)

	for i := range edges {
		edges[i] = g.none
	}

	g.edges = append(g.edges, edges)

	return nil
}

func (g *Graph) SetNode(key string, n int) error {
	idx, present := g.index[key]

	if !present {
		return fmt.Errorf("Node with key %v does not exist", key)
	}

	g.nodes[idx] = n

	return nil
}

func (g *Graph) AddEdge(a, b string, e int) error {
	return g.SetEdge(a, b, e)
}

func (g *Graph) SetEdge(a, b string, e int) error {
	ida, presenta := g.index[a]
	idb, presentb := g.index[b]

	if !presenta {
		return fmt.Errorf("Source node %v does not exist", a)
	}

	if !presentb {
		return fmt.Errorf("Target node %v does not exist", b)
	}

	g.edges[ida][idb] = e

	return nil
}

func (g *Graph) AddEdgeUndirected(a, b string, e int) error {
	return g.SetEdgeUndirected(a, b, e)
}

func (g *Graph) SetEdgeUndirected(a, b string, e int) error {
	err0 := g.SetEdge(a, b, e)

	if err0 != nil {
		return err0
	}

	err1 := g.SetEdge(b, a, e)

	return err1
}

func (g *Graph) Node(key string) int {
	idx, _ := g.index[key]
	return g.nodes[idx]
}

func (g *Graph) Edge(a, b string) int {
	ida, _ := g.index[a]
	idb, _ := g.index[b]
	return g.edges[ida][idb]
}

func (g *Graph) FilterNodes(f func(string, int) bool) []string {
	keys := make([]string, 0)

	for key, idx := range g.index {
		if f(key, g.nodes[idx]) {
			keys = append(keys, key)
		}
	}

	return keys
}

func (g *Graph) Format() string {
	var b strings.Builder

	keys := g.Keys()

	sort.StringSlice(keys).Sort()

	fmt.Fprintf(&b, "\t")

	for _, k := range keys {
		fmt.Fprintf(&b, "%v\t", k)
	}

	fmt.Fprintf(&b, "\n")

	for _, k0 := range keys {
		fmt.Fprintf(&b, "%v\t", k0)

		for _, k1 := range keys {
			e := g.Edge(k0, k1)

			if e == g.none {
				fmt.Fprintf(&b, "%v\t", "-")
			} else {
				fmt.Fprintf(&b, "%v\t", e)
			}
		}

		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}

func (g *Graph) Clone() *Graph {
	size := g.Size()

	index := make(map[string]int, size)

	for k, v := range g.index {
		index[k] = v
	}

	nodes := make([]int, size)
	copy(nodes, g.nodes)

	edges := make([][]int, size)

	for i, es := range g.edges {
		copied := make([]int, len(es))
		copy(copied, es)
		edges[i] = copied
	}

	return &Graph{index, nodes, edges, g.none}
}

func (g *Graph) AddShortestPaths() {
	dist := g.edges

	for k := range g.nodes {
		dist[k][k] = 0
	}

	for k := range g.nodes {
		for i := range g.nodes {
			for j := range g.nodes {
				ik := dist[i][k]
				kj := dist[k][j]

				if ik == g.none || kj == g.none {
					continue
				}

				if dist[i][j] > ik+kj {
					dist[i][j] = ik + kj
				}
			}
		}
	}
}

func parse(file *os.File) (*Graph, error) {
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile("\\AValve ([A-Z]+).+rate=(\\d+);.+ to valves? ([A-Z]+(?:,\\s*[A-Z]+)*)\\z")

	graph := NewGraph()

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindAllStringSubmatch(line, -1)[0]

		key := match[1]
		rate, err := strconv.ParseInt(match[2], 10, 32)

		if err != nil {
			return nil, err
		}

		if graph.HasNode(key) {
			if err := graph.SetNode(key, int(rate)); err != nil {
				return nil, err
			}
		} else {
			if err := graph.AddNode(key, int(rate)); err != nil {
				return nil, err
			}
		}

		neighbors := strings.Split(match[3], ", ")

		for _, nkey := range neighbors {
			if !graph.HasNode(nkey) {
				graph.AddNode(nkey, 0)
			}

			if err := graph.AddEdgeUndirected(key, nkey, 1); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return graph, nil
}

func main() {
	// start := "AA"

	graph, err := parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// 1. Get list of nodes with non-zero flow rate
	// 2. Compute shortest paths between those nodes
	// 3. With budget of 30 mins, visit non-zero flow nodes

	// nonzero := graph.FilterNodes(func(k string, n int) bool {
	// 	return n > 0
	// })

	fmt.Print(graph.Format())

	graph = graph.Clone()
	graph.AddShortestPaths()

	fmt.Print(graph.Format())

	// distances := graph.ShortestPathDistances(0, math.MaxInt, graph.Edge)
}
