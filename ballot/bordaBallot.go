package ballot

import (
	"fmt"
	"net/http"
	"td5/comsoc"
	"time"
)

type BordaBallot struct {
	id       string
	deadline time.Time
	voterIds []string
	votedIds []string
	nbAlts   int
	bordaSwf func(comsoc.Profile) (alts []comsoc.Alternative, err error)
	bordaScf func(comsoc.Profile) (comsoc.Alternative, error)
	profile  comsoc.Profile
}

func NewBordaBallot(
	id string,
	deadline time.Time,
	voterIds []string,
	nbAlts int,
	tieBreakRule []comsoc.Alternative,
) *BordaBallot {
	return &BordaBallot{
		id:       id,
		deadline: deadline,
		voterIds: voterIds,
		nbAlts:   nbAlts,
		bordaSwf: comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(tieBreakRule)),
		bordaScf: comsoc.SCFFactory(comsoc.BordaSCF, comsoc.TieBreakFactory(tieBreakRule)),
	}
}

func (b *BordaBallot) GetId() string {
	return b.id
}

func (b *BordaBallot) GetDeadline() time.Time {
	return b.deadline
}

func (b *BordaBallot) GetWinner() (comsoc.Alternative, error) {
	return b.bordaScf(b.profile)
}

func (b *BordaBallot) GetRanking() ([]comsoc.Alternative, error) {
	return b.bordaSwf(b.profile)
}

func (b *BordaBallot) AddVote(agentId string, vote []comsoc.Alternative, _ []int) (int, error) {
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

func (b *BordaBallot) IsAllowedToVote(agentId string) bool {
	for _, elem := range b.voterIds {
		if elem == agentId {
			return true
		}
	}
	return false
}

func (b *BordaBallot) HasAlreadyVoted(agentId string) bool {
	for _, elem := range b.votedIds {
		if elem == agentId {
			return true
		}
	}
	return false
}
