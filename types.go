package td5

import (
	"time"

	"github.com/StutenEXE/ai30-vote-server/comsoc"
)

type BallotRequest struct {
	Rule         string               `json:"rule"`
	Deadline     time.Time            `json:"deadline"`
	VoterIds     []string             `json:"voter-ids"`
	NbAlts       int                  `json:"#alts"`
	TieBreakRule []comsoc.Alternative `json:"tie-break"`
}

type BallotResponse struct {
	ID string `json:"ballod-id"`
}

type VoteRequest struct {
	AgentId  string               `json:"agent-id"`
	BallotId string               `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  []int                `json:"options"` // facultatif exemple seuil d'acceptation en approval
}

type ResultRequest struct {
	BallotId string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking"`
}
