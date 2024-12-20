package voteserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	td5 "github.com/StutenEXE/ai30-vote-server"

	"github.com/StutenEXE/ai30-vote-server/ballot"
)

type RestVoteServerAgent struct {
	sync.Mutex
	id      string
	port    string
	ballots []ballot.Ballot
}

func NewRestVoteServerAgent(port string) *RestVoteServerAgent {
	return &RestVoteServerAgent{id: port, port: port, ballots: make([]ballot.Ballot, 0)}
}

// Test de la méthode
func (*RestVoteServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestVoteServerAgent) decodeBallotRequest(r *http.Request) (req td5.BallotRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestVoteServerAgent) decodeVoteRequest(r *http.Request) (req td5.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestVoteServerAgent) decodeResultRequest(r *http.Request) (req td5.ResultRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rvsa *RestVoteServerAgent) addBallot(w http.ResponseWriter, r *http.Request) {
	if !rvsa.checkMethod("POST", w, r) {
		return
	}

	req, err := rvsa.decodeBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// On verrouille l'accès aux ballots
	rvsa.Lock()
	defer rvsa.Unlock()
	// création du ballot
	var b ballot.Ballot
	id := fmt.Sprintf("scrutin%d", len(rvsa.ballots)+1)
	if req.Deadline.IsZero() || req.VoterIds == nil || req.NbAlts == 0 || req.TieBreakRule == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("erreur dans la requête, paramètre : rule Deadlin voter-ids #alts tie-break"))
		return
	}

	switch req.Rule {
	case "majority":
		b = ballot.NewMajorityBallot(
			id,
			req.Deadline,
			req.VoterIds,
			req.NbAlts,
			req.TieBreakRule,
		)
	case "borda":
		b = ballot.NewBordaBallot(
			id,
			req.Deadline,
			req.VoterIds,
			req.NbAlts,
			req.TieBreakRule,
		)
	case "approval":
		b = ballot.NewApprovalBallot(
			id,
			req.Deadline,
			req.VoterIds,
			req.NbAlts,
			req.TieBreakRule,
		)
	case "copeland":
		b = ballot.NewCopelandBallot(
			id,
			req.Deadline,
			req.VoterIds,
			req.NbAlts,
			req.TieBreakRule,
		)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("règle inconnue '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}

	rvsa.ballots = append(rvsa.ballots, b)
	w.WriteHeader(http.StatusCreated)
	response := td5.BallotResponse{
		ID: id,
	}
	byt, _ := json.Marshal(response)
	w.Write(byt)
}

func (rvsa *RestVoteServerAgent) getBallotById(id string) ballot.Ballot {
	for _, b := range rvsa.ballots {
		if b.GetId() == id {
			return b
		}
	}
	return nil
}

func (rvsa *RestVoteServerAgent) addVote(w http.ResponseWriter, r *http.Request) {
	if !rvsa.checkMethod("POST", w, r) {
		return
	}

	req, err := rvsa.decodeVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// On verrouille l'accès aux ballots
	rvsa.Lock()
	defer rvsa.Unlock()
	bal := rvsa.getBallotById(req.BallotId)
	if bal == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "pas de ballot avec l'id %v", req.BallotId)
		return
	}
	// Deadline dépassée
	if bal.GetDeadline().Before(time.Now()) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "la deadline est dépassée")
		return
	}
	// A déjà voté
	if bal.HasAlreadyVoted(req.AgentId) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "vote déjà effectué")
		return
	}
	// Pas le droit de vote
	if !bal.IsAllowedToVote(req.AgentId) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "pas autorisé à voter ")
		return
	}

	code, err := bal.AddVote(req.AgentId, req.Prefs, req.Options)

	if err != nil {
		w.WriteHeader(code)
		fmt.Fprint(w, err.Error())
		return
	}
	//vote effectué
	w.WriteHeader(code)
}

func (rvsa *RestVoteServerAgent) giveResult(w http.ResponseWriter, r *http.Request) {
	if !rvsa.checkMethod("POST", w, r) {
		return
	}
	req, err := rvsa.decodeResultRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	bal := rvsa.getBallotById(req.BallotId)
	if bal == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "pas de ballot avec l'id %v", req.BallotId)
		return
	}
	//Deadline pas encore passée
	if bal.GetDeadline().After(time.Now()) {
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprintf(w, "la deadline n'est pas encore passée")
		return
	}

	winner, err := bal.GetWinner()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	ranking, err := bal.GetRanking()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	resp := td5.ResultResponse{Winner: winner, Ranking: ranking}
	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.Write(serial)
}

func (rvsa *RestVoteServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rvsa.addBallot)
	mux.HandleFunc("/vote", rvsa.addVote)
	mux.HandleFunc("/result", rvsa.giveResult)

	// création du serveur http
	s := &http.Server{
		Addr:           rvsa.port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rvsa.port)
	go log.Fatal(s.ListenAndServe())
}
