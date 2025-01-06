package main

import "sync"

// Contains unique identifiers(ID) and the network address
type Node struct {
	ID   string
	Addr string
}

// Hashes for nodes with concurrency protection
type HashRing struct {
	nodes  []Node
	hashes []uint32
	lock   sync.RWMutex
}

func NewHashRing() *HashRing {
	return &HashRing{}
}

func (h *HashRing) AddNode(node Node) {
	return
}

func (h *HashRing) RemoveNode(nodeID string) {
	return
}

func (h *HashRing) GetNode(key string) Node {
	return Node{}
}

func (h *HashRing) hash(key string) uint32 {
	return 0
}
