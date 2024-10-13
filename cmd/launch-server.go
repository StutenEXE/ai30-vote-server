package main

import (
	"fmt"
	"time"

	td5 "github.com/StutenEXE/ai30-vote-server"
	"github.com/StutenEXE/ai30-vote-server/agent"
	"github.com/StutenEXE/ai30-vote-server/comsoc"
	"github.com/StutenEXE/ai30-vote-server/voteserveragent"
)

func main() {
	// Lancement du serveur de vote
	server := voteserveragent.NewRestVoteServerAgent(":8080")
	go server.Start()
	defer fmt.Scanln()

	// Création d'un ballot de type majority

	ballotReq := td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(10 * time.Second),
		VoterIds:     []string{"Clément", "Alexandre"},
		NbAlts:       2,
		TieBreakRule: []comsoc.Alternative{1, 2},
	}

	//créée le ballot
	_, ballotID, err := agent.Ballot(ballotReq)

	if err != nil {
		fmt.Println("Erreur lors de la création du ballot:", err)
		return
	}

	// vote 1

	voteReq1 := td5.VoteRequest{
		AgentId:  "Clément",
		BallotId: ballotID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
	}

	_, err = agent.Vote(voteReq1)

	if err != nil {
		fmt.Println("Erreur lors du vote de l'agent1 1:", err)
	} else {

		fmt.Println("Vote de l'agent1 1 enregistré avec succès.")
	}

	//vote 2

	voteReq2 := td5.VoteRequest{
		AgentId:  "Alexandre",
		BallotId: ballotID.ID,
		Prefs:    []comsoc.Alternative{2, 1},
	}

	_, err = agent.Vote(voteReq2)

	if err != nil {
		fmt.Println("Erreur lors du vote de l'agent2 2:", err)
	} else {

		fmt.Println("Vote de l'agent2 enregistré avec succès.")
	}

	// On attend que la deadline soit dépassée pour demander le résultat
	time.Sleep(10 * time.Second)

	//  demande le résultat

	resultReq := td5.ResultRequest{
		BallotId: ballotID.ID,
	}

	_, result, err := agent.Result(resultReq)

	if err != nil {
		fmt.Println("Erreur lors de la récupération des résultats:", err)
	} else {
		fmt.Printf("Résultat du scrutin: Gagnant - %v, Classement - %v\n", result.Winner, result.Ranking)
	}
}
