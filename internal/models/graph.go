package models

import "fmt"

type Node struct {
	Id string
}

type Edge struct {
	From string
	To   string
	Name string
}

type Graph struct {
	start Node
	nodes map[string]Node
	edges map[string]Edge
}

func CreateGraph() Graph {
	return Graph{
		nodes: map[string]Node{},
		edges: map[string]Edge{},
	}
}

func (g *Graph) AddNode(id string) {
	node := Node{Id: id}
	if len(g.nodes) == 0 {
		g.start = node
	}
	g.nodes[id] = node
}

func edgeKey(from, name string) string {
	key := fmt.Sprintf("%s-%s", from, name)
	return key
}

func (g *Graph) AddEdge(from, to, name string) {
	if _, ok := g.nodes[from]; !ok {
		return
	}
	if _, ok := g.nodes[to]; !ok {
		return
	}
	key := edgeKey(from, name)
	g.edges[key] = Edge{
		From: from,
		To:   to,
		Name: name,
	}
}

func (g *Graph) GetEdge(from, name string) *Edge {
	key := edgeKey(from, name)
	if edge, ok := g.edges[key]; ok {
		return &edge
	}
	return nil
}

func (g *Graph) GetStart() *Node {
	return &g.start
}
