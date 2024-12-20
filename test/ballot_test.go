package test

import (
	"testing"
	"time"

	td5 "github.com/StutenEXE/ai30-vote-server"

	"github.com/StutenEXE/ai30-vote-server/comsoc"
	"github.com/StutenEXE/ai30-vote-server/voteclientagent"
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

	// Règle inconnue
	status, ballotID, _ := voteclientagent.Ballot(ballotReq)
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
	status, ballotID, _ = voteclientagent.Ballot(ballotReq)
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
	status, ballotID, _ = voteclientagent.Ballot(ballotReq)
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
	status, ballotID, _ = voteclientagent.Ballot(ballotReq)
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
	status, ballotID, _ = voteclientagent.Ballot(ballotReq)
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
	status, ballotID, _ = voteclientagent.Ballot(ballotReq)
	if ballotID.ID == "" || status != "201 Created" {
		t.Error(status)
	}
}
