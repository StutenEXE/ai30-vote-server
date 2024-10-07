package ballot

import (
	"fmt"
	"net/http"
	"time"

	"td5/comsoc"
)

type MajorityBallot struct {
	id           string
	deadline     time.Time
	voterIds     []string
	votedIds     []string
	nbAlts       int
	tieBreakRule []int
	profile      comsoc.Profile
}

func (b *MajorityBallot) GetId() string {
	return b.id
}

func (b *MajorityBallot) GetDeadline() time.Time {
	return b.deadline
}

func (b *MajorityBallot) GetWinner() (int, error) {
	// TODO
	return 0, nil
}

func (b *MajorityBallot) GetRanking() ([]int, error) {
	// TODO
	return []int{0}, nil
}

func (b *MajorityBallot) AddVote(agentId string, vote []int, options []int) (int, error) {
	// Deadline dépassée
	if b.GetDeadline().Before(time.Now()) {
		return http.StatusServiceUnavailable, fmt.Errorf("la deadline est dépassée")
	}
	// A déjà voté
	if b.hasAlreadyVoted(agentId) {
		return http.StatusForbidden, fmt.Errorf("vote déjà effectué", b.deadline)
	}
	// Pas le droit de vote
	if b.allowedToVote(agentId) == false {
		return http.StatusBadRequest, fmt.Errorf("pas autorisé à voter")
	}

	return http.StatusOK, nil
}

func (b *MajorityBallot) allowedToVote(agentId string) bool {
	for _, elem := range b.voterIds {
		if elem == agentId {
			return true
		}
	}
	return false
}

func (b *MajorityBallot) hasAlreadyVoted(agentId string) bool {
	for _, elem := range b.votedIds {
		if elem == agentId {
			return true
		}
	}
	return false
}
