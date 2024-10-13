package test

import (
	"reflect"
	"testing"
	"time"

	td5 "github.com/StutenEXE/ai30-vote-server"

	"github.com/StutenEXE/ai30-vote-server/agent"
	"github.com/StutenEXE/ai30-vote-server/comsoc"
)

func TestDeadlineNotPass(t *testing.T) {
	// on crée le ballot
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

func TestMajority(t *testing.T) {
	listAgt := []string{"ag1", "ag2", "ag3", "ag4", "ag5"}
	listAlt := []comsoc.Alternative{1, 2, 3, 4}
	// On crée un ballot de 5 candidats
	ballotReq := td5.BallotRequest{
		Rule:         "majority",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     listAgt,
		NbAlts:       4,
		TieBreakRule: listAlt,
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// On crée 5 profils de votes tels que 1 et 3 soient à égalité
	// D'après la règle de tie break 1 devrait être élu
	listVotingProfiles := comsoc.Profile{
		[]comsoc.Alternative{1, 3, 2, 4},
		[]comsoc.Alternative{3, 1, 2, 4},
		[]comsoc.Alternative{1, 3, 4, 2},
		[]comsoc.Alternative{3, 1, 4, 2},
		[]comsoc.Alternative{2, 4, 3, 1},
	}
	for i, prof := range listVotingProfiles {
		voteReq := td5.VoteRequest{
			AgentId:  listAgt[i],
			BallotId: bID.ID,
			Prefs:    prof,
		}
		status, _ = agent.Vote(voteReq)
		if status != "200 OK" {
			t.Error(status)
		}
	}
	// On attends la fin de la deadline
	time.Sleep(time.Second * 2)

	// On demande le résultat
	resultReq := td5.ResultRequest{
		BallotId: bID.ID,
	}

	status, result, _ := agent.Result(resultReq)
	if status != "200 OK" {
		t.Error(status)
	}
	if result.Winner != 1 {
		t.Errorf("Majority : winner should be 1 but is %d", result.Winner)
	}
	expectedRanking := []comsoc.Alternative{1, 3}
	if !reflect.DeepEqual(result.Ranking, expectedRanking) {
		t.Errorf("Majority : ranking should be %v but is %v", expectedRanking, result.Ranking)
	}
}

func TestApproval(t *testing.T) {
	listAgt := []string{"ag1", "ag2", "ag3", "ag4", "ag5"}
	listAlt := []comsoc.Alternative{1, 2, 3, 4}
	// On crée un ballot de 5 candidats
	ballotReq := td5.BallotRequest{
		Rule:         "approval",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     listAgt,
		NbAlts:       4,
		TieBreakRule: listAlt,
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// On crée 5 profils de votes tels que 1 et 3 soient à égalité
	// D'après la règle de tie break 1 devrait être élu
	listVotingProfiles := comsoc.Profile{
		[]comsoc.Alternative{1, 3, 2, 4},
		[]comsoc.Alternative{1, 2, 3, 4},
		[]comsoc.Alternative{1, 3, 4, 2},
		[]comsoc.Alternative{2, 3, 1, 4},
		[]comsoc.Alternative{1, 3, 4, 2},
	}
	listThresholds := []int{2, 3, 2, 3, 2}
	for i, prof := range listVotingProfiles {
		voteReq := td5.VoteRequest{
			AgentId:  listAgt[i],
			BallotId: bID.ID,
			Prefs:    prof,
			Options:  []int{listThresholds[i]},
		}
		status, _ = agent.Vote(voteReq)
		if status != "200 OK" {
			t.Error(status)
		}
	}
	// On attends la fin de la deadline
	time.Sleep(time.Second * 2)

	// On demande le résultat
	resultReq := td5.ResultRequest{
		BallotId: bID.ID,
	}

	status, result, _ := agent.Result(resultReq)
	if status != "200 OK" {
		t.Error(status)
	}
	if result.Winner != 1 {
		t.Errorf("Approval : winner should be 1 but is %d", result.Winner)
	}
	expectedRanking := []comsoc.Alternative{1, 3}
	if !reflect.DeepEqual(result.Ranking, expectedRanking) {
		t.Errorf("Approval : ranking should be %v but is %v", expectedRanking, result.Ranking)
	}
}

func TestBorda(t *testing.T) {
	listAgt := []string{"ag1", "ag2", "ag3", "ag4", "ag5"}
	listAlt := []comsoc.Alternative{1, 2, 3, 4}
	// On crée un ballot de 5 candidats
	ballotReq := td5.BallotRequest{
		Rule:         "borda",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     listAgt,
		NbAlts:       4,
		TieBreakRule: listAlt,
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// On crée 5 profils de votes tels que 1 et 3 soient à égalité
	// D'après la règle de tie break 1 devrait être élu
	listVotingProfiles := comsoc.Profile{
		[]comsoc.Alternative{1, 3, 2, 4},
		[]comsoc.Alternative{3, 1, 2, 4},
		[]comsoc.Alternative{1, 3, 4, 2},
		[]comsoc.Alternative{3, 1, 4, 2},
		[]comsoc.Alternative{2, 1, 3, 4},
	}
	for i, prof := range listVotingProfiles {
		voteReq := td5.VoteRequest{
			AgentId:  listAgt[i],
			BallotId: bID.ID,
			Prefs:    prof,
		}
		status, _ = agent.Vote(voteReq)
		if status != "200 OK" {
			t.Error(status)
		}
	}
	// On attends la fin de la deadline
	time.Sleep(time.Second * 2)

	// On demande le résultat
	resultReq := td5.ResultRequest{
		BallotId: bID.ID,
	}

	status, result, _ := agent.Result(resultReq)
	if status != "200 OK" {
		t.Error(status)
	}
	if result.Winner != 1 {
		t.Errorf("Borda : winner should be 1 but is %d", result.Winner)
	}
	expectedRanking := []comsoc.Alternative{1, 3}
	if !reflect.DeepEqual(result.Ranking, expectedRanking) {
		t.Errorf("Borda : ranking should be %v but is %v", expectedRanking, result.Ranking)
	}
}

func TestCopeland(t *testing.T) {
	listAgt := []string{"ag1", "ag2", "ag3", "ag4", "ag5"}
	listAlt := []comsoc.Alternative{1, 2, 3, 4}
	// On crée un ballot de 5 candidats
	ballotReq := td5.BallotRequest{
		Rule:         "borda",
		Deadline:     time.Now().Add(2 * time.Second),
		VoterIds:     listAgt,
		NbAlts:       4,
		TieBreakRule: listAlt,
	}
	status, bID, _ := agent.Ballot(ballotReq)
	if bID.ID == "" || status != "201 Created" {
		t.Error(status)
	}

	// On crée 5 profils de votes tels que 1 et 3 soient à égalité
	// D'après la règle de tie break 1 devrait être élu
	listVotingProfiles := comsoc.Profile{
		[]comsoc.Alternative{1, 3, 2, 4},
		[]comsoc.Alternative{3, 1, 2, 4},
		[]comsoc.Alternative{1, 3, 4, 2},
		[]comsoc.Alternative{3, 1, 4, 2},
		[]comsoc.Alternative{2, 1, 3, 4},
	}
	for i, prof := range listVotingProfiles {
		voteReq := td5.VoteRequest{
			AgentId:  listAgt[i],
			BallotId: bID.ID,
			Prefs:    prof,
		}
		status, _ = agent.Vote(voteReq)
		if status != "200 OK" {
			t.Error(status)
		}
	}
	// On attends la fin de la deadline
	time.Sleep(time.Second * 2)

	// On demande le résultat
	resultReq := td5.ResultRequest{
		BallotId: bID.ID,
	}

	status, result, _ := agent.Result(resultReq)
	if status != "200 OK" {
		t.Error(status)
	}
	if result.Winner != 1 {
		t.Errorf("Copeland : winner should be 1 but is %d", result.Winner)
	}
	expectedRanking := []comsoc.Alternative{1, 3}
	if !reflect.DeepEqual(result.Ranking, expectedRanking) {
		t.Errorf("Copeland : ranking should be %v but is %v", expectedRanking, result.Ranking)
	}
}
