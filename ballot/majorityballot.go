package ballot

import (
	"fmt"
	"net/http"
	"time"

	"td5/comsoc"
)

type MajorityBallot struct {
	id          string
	deadline    time.Time
	voterIds    []string
	votedIds    []string
	nbAlts      int
	majoritySwf func(comsoc.Profile) (alts []comsoc.Alternative, err error)
	majorityScf func(comsoc.Profile) (comsoc.Alternative, error)
	profile     comsoc.Profile
}

func NewMajorityBallot(
	id string,
	deadline time.Time,
	voterIds []string,
	nbAlts int,
	tieBreakRule []comsoc.Alternative,
) *MajorityBallot {
	return &MajorityBallot{
		id:          id,
		deadline:    deadline,
		voterIds:    voterIds,
		nbAlts:      nbAlts,
		majoritySwf: comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(tieBreakRule)),
		majorityScf: comsoc.SCFFactory(comsoc.MajoritySCF, comsoc.TieBreakFactory(tieBreakRule)),
	}
}

func (b *MajorityBallot) GetId() string {
	return b.id
}

func (b *MajorityBallot) GetDeadline() time.Time {
	return b.deadline
}

func (b *MajorityBallot) GetWinner() (comsoc.Alternative, error) {
	return b.majorityScf(b.profile)
}

func (b *MajorityBallot) GetRanking() ([]comsoc.Alternative, error) {
	return b.majoritySwf(b.profile)
}

func (b *MajorityBallot) AddVote(agentId string, vote []comsoc.Alternative, _ []int) (int, error) {
	// Vote invalide - pas le bon nombre d'alternatives
	err := CheckAlternativesUnicity(vote)
	if len(vote) != b.nbAlts || err != nil {
		return http.StatusBadRequest, fmt.Errorf("vote invalide ")
	}

	// Ajout du vote
	b.votedIds = append(b.votedIds, agentId)
	b.profile = append(b.profile, vote)

	return http.StatusOK, nil
}

func (b *MajorityBallot) IsAllowedToVote(agentId string) bool {
	for _, elem := range b.voterIds {
		if elem == agentId {
			return true
		}
	}
	return false
}

func (b *MajorityBallot) HasAlreadyVoted(agentId string) bool {
	for _, elem := range b.votedIds {
		if elem == agentId {
			return true
		}
	}
	return false
}
