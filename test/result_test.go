package test

import (
	"td5"
	"td5/agent"
	"td5/comsoc"
	"testing"
	"time"
)

func TestDeadLineNotPass(t *testing.T) {
	// on crée le ballot avant du coup
	ballotReq := td5.BallotRequest{
		Rule:         "approval",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     []string{"agent1", "agent2"},
		NbAlts:       2,
		TieBreakRule: []comsoc.Alternative{1, 2},
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// ON vérifie le ballot avec une bonne requête
	voteReq := td5.VoteRequest{
		AgentId:  "agent1",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
		Options:  []int{1},
	}
	status, _ = agent.Vote(voteReq)
	if status != "200 OK" {
		t.Error(status)
	}
	// 2ème vote pour le résultat
	voteReq = td5.VoteRequest{
		AgentId:  "agent2",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{2, 1},
		Options:  []int{2},
	}
	status, _ = agent.Vote(voteReq)
	if status != "200 OK" {
		t.Error(status)
	}

	resultReq := td5.ResultRequest{
		BallotId: bID.ID,
	}

	//On demande le résultat trop tot
	status, result, _ := agent.Result(resultReq)
	if status != "425 Too Early" {
		t.Error(status)
	}
	time.Sleep(time.Second * 2)
	status, result, _ = agent.Result(resultReq)

	if status != "200 OK" || result.Winner != 1 {
		t.Error(status)
	}
}
