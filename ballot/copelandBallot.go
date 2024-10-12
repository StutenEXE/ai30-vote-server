package ballot

import (
	"fmt"
	"net/http"
	"time"

	"td5/comsoc"
)

type CopelandBallot struct {
	id          string
	deadline    time.Time
	voterIds    []string
	votedIds    []string
	nbAlts      int
	copelandSwf func(comsoc.Profile) (alts []comsoc.Alternative, err error)
	copelandScf func(comsoc.Profile) (comsoc.Alternative, error)
	profile     comsoc.Profile
}

func NewCopelandBallot(
	id string,
	deadline time.Time,
	voterIds []string,
	nbAlts int,
	tieBreakRule []comsoc.Alternative,
) *CopelandBallot {
	return &CopelandBallot{
		id:          id,
		deadline:    deadline,
		voterIds:    voterIds,
		nbAlts:      nbAlts,
		copelandSwf: comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(tieBreakRule)),
		copelandScf: comsoc.SCFFactory(comsoc.CopelandSCF, comsoc.TieBreakFactory(tieBreakRule)),
	}
}

func (b *CopelandBallot) GetId() string {
	return b.id
}

func (b *CopelandBallot) GetDeadline() time.Time {
	return b.deadline
}

func (b *CopelandBallot) GetWinner() (comsoc.Alternative, error) {
	return b.copelandScf(b.profile)
}

func (b *CopelandBallot) GetRanking() ([]comsoc.Alternative, error) {
	return b.copelandSwf(b.profile)
}

func (b *CopelandBallot) AddVote(agentId string, vote []comsoc.Alternative, _ []int) (int, error) {
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

func (b *CopelandBallot) IsAllowedToVote(agentId string) bool {
	for _, elem := range b.voterIds {
		if elem == agentId {
			return true
		}
	}
	return false
}

func (b *CopelandBallot) HasAlreadyVoted(agentId string) bool {
	for _, elem := range b.votedIds {
		if elem == agentId {
			return true
		}
	}
	return false
}
