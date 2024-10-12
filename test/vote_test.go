package test

import (
	"td5"
	"td5/agent"
	"td5/comsoc"
	"testing"
	"time"
)

func TestVoteWrongParameter(t *testing.T) {

	voteReq := td5.VoteRequest{
		AgentId:  "agent1",
		BallotId: "test",
		Prefs:    []comsoc.Alternative{1, 2},
	}

	// Pas de bon ballot id
	status, _ := agent.Vote(voteReq)
	if status != "400 Bad Request" {
		t.Error(status)
	}

	// On crée le ballot avant du coup
	ballotReq := td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     []string{"agent1", "agent2"},
		NbAlts:       2,
		TieBreakRule: []comsoc.Alternative{1, 2},
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// On vérifie le ballot avec une bonne requête
	voteReq = td5.VoteRequest{
		AgentId:  "agent1",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
	}
	status, _ = agent.Vote(voteReq)
	if status != "200 OK" {
		t.Error(status)
	}

	// A déjà voté
	status, _ = agent.Vote(voteReq)
	if status != "403 Forbidden" {
		t.Error(status)
	}

	// Agent inconnu
	voteReq = td5.VoteRequest{
		AgentId:  "Clement",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
	}
	// A déjà voté
	status, _ = agent.Vote(voteReq)
	if status != "400 Bad Request" {
		t.Error(status)
	}

	voteReq = td5.VoteRequest{
		AgentId:  "Clement",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{4, 2},
	}
	// Mauvais vote
	status, _ = agent.Vote(voteReq)
	if status != "400 Bad Request" {
		t.Error(status)
	}

	voteReq = td5.VoteRequest{
		AgentId:  "agent2",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
	}
	time.Sleep(time.Second * 2)
	// Deadline
	status, _ = agent.Vote(voteReq)
	if status != "503 Service Unavailable" {
		t.Error(status)
	}

}

func TestVoteApprovalWrongOptions(t *testing.T) {
	// On crée le ballot
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

	// On vérifie le ballot avec une bonne requête
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

	// Pas d'options
	voteReq = td5.VoteRequest{
		AgentId:  "agent2",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
	}
	status, _ = agent.Vote(voteReq)
	if status != "400 Bad Request" {
		t.Error(status)
	}

	// Options pas bonne
	voteReq = td5.VoteRequest{
		AgentId:  "agent2",
		BallotId: bID.ID,
		Prefs:    []comsoc.Alternative{1, 2},
		Options:  []int{77},
	}
	status, _ = agent.Vote(voteReq)
	if status != "400 Bad Request" {
		t.Error(status)
	}
}
