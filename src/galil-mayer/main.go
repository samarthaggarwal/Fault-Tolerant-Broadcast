package main

import (
	"fmt"
	"src/node"
)

func main() {
	fmt.Println("hello world!")

	n := node.Node{id:1234}
	fmt.Println("id = {}", n.get_id())
}
