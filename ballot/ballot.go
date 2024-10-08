package ballot

import (
	"sort"
	"td5/comsoc"
	"time"
)

type Ballot interface {
	GetId() string
	GetDeadline() time.Time
	GetWinner() (comsoc.Alternative, error)
	GetRanking() ([]comsoc.Alternative, error)
	AddVote(agentId string, vote []comsoc.Alternative, options []int) (int, error)
	HasAlreadyVoted(agentId string) bool
	IsAllowedToVote(agentId string) bool
}

func getRankingFromCount(c comsoc.Count) []comsoc.Alternative {
	alts := make([]comsoc.Alternative, 0, len(c))

	for alt := range c {
		alts = append(alts, alt)
	}

	sort.SliceStable(alts, func(i, j int) bool {
		return c[alts[i]] > c[alts[j]]
	})

	return alts
}
