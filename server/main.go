// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/getamis/alice/server/peer"
)

func main() {
	var wg sync.WaitGroup
	var nodecount int = 3
	var initiatingNode *peer.Node
	peermap := make(map[string]*peer.PeerInfo)

	for i := 0; i < nodecount; i++ {
		wg.Add(1)
		var nodeid string = "node-" + strconv.Itoa(i)
		node := peer.NewNode(nodeid, (uint32(8080 + i)))
		if i == 0 {
			initiatingNode = node
		} else {
			peermap[node.ID()] = peer.NewPeerInfo(nodeid, "localhost", uint32(8080+i), 0)
		}
		go func(node *peer.Node) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			node.Start()
		}(node)
	}
	go func(initiatingNode *peer.Node) {
		log.Println("initiaiting...")
		defer wg.Done()
		initiatingNode.StartDKGSession(peermap, 2, 0, "")
	}(initiatingNode)
	wg.Wait()
	log.Println("main")
}

func generateKeys(node *peer.Node) {
	log.Println("node...")
}
