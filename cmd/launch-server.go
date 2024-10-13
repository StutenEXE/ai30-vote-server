package main

import (
	"fmt"

	"github.com/StutenEXE/ai30-vote-server/voteserveragent"
)

func main() {
	// Lancement du serveur de vote
	server := voteserveragent.NewRestVoteServerAgent(":8080")
	go server.Start()
	defer fmt.Scanln()
}
