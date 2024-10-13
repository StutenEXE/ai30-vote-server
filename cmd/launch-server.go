package main

import (
	"fmt"

	"github.com/StutenEXE/ai30-vote-server/voteserveragent"
)

func main() {
	server := voteserveragent.NewRestVoteServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
