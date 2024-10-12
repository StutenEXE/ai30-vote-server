package test

import (
	"td5"
	"td5/agent"
	"td5/comsoc"
	"testing"
	"time"
)

func TestWrongBallotParameters(t *testing.T) {

	// Création d'un ballot (peu importe vu que on teste les mauvais paramètre)
	ballotReq := td5.BallotRequest{
		Rule:         "LePlusBeau",
		Deadline:     time.Now().Add(1 * time.Hour),
		VoterIds:     []string{"agent1", "agent2"},
		NbAlts:       3,
		TieBreakRule: []comsoc.Alternative{},
	}

	// 1 Rule inconnu
	status, ballotID, _ := agent.Ballot(ballotReq)
	if ballotID.ID != "" || status != "501 Not Implemented" {
		t.Error(status)
	}

	ballotReq = td5.BallotRequest{
		Rule:         "majority",
		VoterIds:     []string{"agent1", "agent2"},
		NbAlts:       3,
		TieBreakRule: []comsoc.Alternative{},
	}

	// Pas assez d'arg (deadline)
	status, ballotID, _ = agent.Ballot(ballotReq)
	if ballotID.ID != "" || status != "400 Bad Request" {
		t.Error(status)
	}

	ballotReq = td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(1 * time.Hour),
		NbAlts:       3,
		TieBreakRule: []comsoc.Alternative{},
	}

	//Pas assez d'arg (voterIds)
	status, ballotID, _ = agent.Ballot(ballotReq)
	if ballotID.ID != "" || status != "400 Bad Request" {
		t.Error(status)
	}

	ballotReq = td5.BallotRequest{
		Rule:         "majority",
		VoterIds:     []string{"agent1", "agent2"},
		Deadline:     time.Now().Add(1 * time.Hour),
		TieBreakRule: []comsoc.Alternative{},
	}

	// Pas assez d'arg (nbAlts)
	status, ballotID, _ = agent.Ballot(ballotReq)
	if ballotID.ID != "" || status != "400 Bad Request" {
		t.Error(status)
	}

	ballotReq = td5.BallotRequest{
		Rule:     "majority",
		VoterIds: []string{"agent1", "agent2"},
		Deadline: time.Now().Add(1 * time.Hour),
		NbAlts:   3,
	}

	// Pas assez d'arg (TieBreakRule)
	status, ballotID, _ = agent.Ballot(ballotReq)
	if ballotID.ID != "" || status != "400 Bad Request" {
		t.Error(status)
	}

	ballotReq = td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(1 * time.Hour),
		VoterIds:     []string{"agent1", "agent2"},
		NbAlts:       3,
		TieBreakRule: []comsoc.Alternative{},
	}

	// Bonne requete
	status, ballotID, _ = agent.Ballot(ballotReq)
	if ballotID.ID == "" || status != "201 Created" {
		t.Error(status)
	}
}
