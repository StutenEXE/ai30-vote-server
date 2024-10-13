package main

import (
	"fmt"
	"time"

	td5 "github.com/StutenEXE/ai30-vote-server"
	"github.com/StutenEXE/ai30-vote-server/comsoc"
	"github.com/StutenEXE/ai30-vote-server/voteserveragent"
)

func main() {
	// Lancement du serveur de vote
	server := voteserveragent.NewRestVoteServerAgent(":8080")
	go server.Start()

	// Création d'un ballot de type majority
	ballotReq := td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(10 * time.Second),
		VoterIds:     []string{"td5", "td52"},
		NbAlts:       2,
		TieBreakRule: []comsoc.Alternative{1, 2},
	}

	//créée le ballot
	ballotID, err := td5.Ballot(ballotReq)
	if err != nil {
		fmt.Println("Erreur lors de la création du ballot:", err)
		return
	}
	fmt.Println("Ballot créé avec succès, ID:", ballotID)

	// vote 1
	voteReq1 := td5.VoteRequest{
		td5Id:    "td5",
		BallotId: ballotID,
		Prefs:    []comsoc.Alternative{1, 2},
	}
	err = td5.Vote(voteReq1)
	if err != nil {
		fmt.Println("Erreur lors du vote de l'agent1 1:", err)
	} else {
		fmt.Println("Vote de l'agent1 1 enregistré avec succès.")
	}

	//vote 2
	voteReq2 := td5.VoteRequest{
		td5Id:    "td52",
		BallotId: ballotID,
		Prefs:    []comsoc.Alternative{2, 1},
	}
	err = td5.Vote(voteReq2)
	if err != nil {
		fmt.Println("Erreur lors du vote de l'agent2 2:", err)
	} else {
		fmt.Println("Vote de l'agent2 enregistré avec succès.")
	}

	// On attend que la deadline soit dépassée pour demander le résultat
	time.Sleep(10 * time.Second)

	//  demande le résultat
	resultReq := td5.ResultRequest{
		BallotId: ballotID,
	}
	result, err := td5.Result(resultReq)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des résultats:", err)
	} else {
		fmt.Printf("Résultat du scrutin: Gagnant - %v, Classement - %v\n", result.Winner, result.Ranking)
	}
	fmt.Scanln()
}
