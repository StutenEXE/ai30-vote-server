package ballot

import (
	"fmt"
	"net/http"
	"time"

	"github.com/StutenEXE/ai30-vote-server/comsoc"
)

type ApprovalBallot struct {
	id           string
	deadline     time.Time
	voterIds     []string
	votedIds     []string
	nbAlts       int
	profile      comsoc.Profile
	thresholds   []int
	approvalScf  func(comsoc.Profile, []int) ([]comsoc.Alternative, error)
	tieBreakRule []comsoc.Alternative
}

func NewApprovalBallot(
	id string,
	deadline time.Time,
	voterIds []string,
	nbAlts int,
	tieBreakRule []comsoc.Alternative,
) *ApprovalBallot {
	return &ApprovalBallot{
		id:           id,
		deadline:     deadline,
		voterIds:     voterIds,
		nbAlts:       nbAlts,
		approvalScf:  comsoc.ApprovalSCF,
		tieBreakRule: tieBreakRule,
	}
}

func (b *ApprovalBallot) GetId() string {
	return b.id
}

func (b *ApprovalBallot) GetDeadline() time.Time {
	return b.deadline
}

func (b *ApprovalBallot) GetWinner() (comsoc.Alternative, error) {
	appro2 := func(p comsoc.Profile) (alt []comsoc.Alternative, err error) {
		return comsoc.ApprovalSCF(p, b.thresholds)
	}
	return comsoc.SCFFactory(appro2, comsoc.TieBreakFactory(b.tieBreakRule))(b.profile)
}

func (b *ApprovalBallot) GetRanking() ([]comsoc.Alternative, error) {
	appro2 := func(p comsoc.Profile) (c comsoc.Count, err error) {
		return comsoc.ApprovalSWF(p, b.thresholds)
	}
	c, err := comsoc.SWFFactory(appro2, comsoc.TieBreakFactory(b.tieBreakRule))(b.profile)
	return c, err
}

func (b *ApprovalBallot) AddVote(agentId string, vote []comsoc.Alternative, options []int) (int, error) {
	// Vote invalide - pas le bon nombre d'alternatives
	err := CheckAlternativesUnicity(vote)
	if len(vote) != b.nbAlts || err != nil {
		return http.StatusBadRequest, fmt.Errorf("vote invalide ")
	}

	if options == nil || len(options) > 1 || options[0] > len(vote) || options[0] < 0 {
		return http.StatusBadRequest, fmt.Errorf("vote invalide")
	}

	// Ajout du vote
	b.votedIds = append(b.votedIds, agentId)
	b.profile = append(b.profile, vote)
	b.thresholds = append(b.thresholds, options[0])
	return http.StatusOK, nil
}

func (b *ApprovalBallot) IsAllowedToVote(agentId string) bool {
	for _, elem := range b.voterIds {
		if elem == agentId {
			return true
		}
	}
	return false
}

func (b *ApprovalBallot) HasAlreadyVoted(agentId string) bool {
	for _, elem := range b.votedIds {
		if elem == agentId {
			return true
		}
	}
	return false
}
