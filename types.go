package td5

import "time"

type BallotRequest struct {
	Rule         string    `json:"rule"`
	Deadline     time.Time `json:"deadline"`
	VoterIds     []string  `json:"voter-ids"`
	NbAlts       int       `json:"nb-alts"`
	TieBreakRule []int     `json:"tie-break"`
}

type VoteRequest struct {
	AgentId  string `json:"agent-id"`
	BallotId string `json:"ballot-id"`
	Prefs    []int  `json:"prefs"`
	Options  []int  `json:"options"` // facultatif exemple seuil d'acceptation en approval
}

type ResultRequest struct {
	BallotId string `json:"scrutin12"`
}

type ResultResponse struct {
	Winner  int   `json:"winner"`
	Ranking []int `json:"ranking"`
}
