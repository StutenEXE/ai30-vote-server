package ballot

import (
	"time"
)

type Ballot interface {
	GetId() string
	GetDeadline() time.Time
	GetWinner() (int, error)
	GetRanking() ([]int, error)
	AddVote(agentId string, vote []int, options []int) (int, error)
}
