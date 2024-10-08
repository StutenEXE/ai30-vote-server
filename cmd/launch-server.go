package main

import (
	"fmt"
	"td5/voteserveragent"
)

func main() {
	server := voteserveragent.NewRestVoteServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
