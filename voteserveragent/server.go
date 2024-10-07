package voteserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"td5"
	"td5/ballot"
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
func (rvsa *RestVoteServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
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

	// création du ballot
	var b *ballot.Ballot
	switch req.Rule {
	case "majority":
		b = ballot.MajorityBallot{
			id:           fmt.FormatString("scrutin%d", len(rvsa.ballots)+1),
			deadline:     req.Deadline,
			voterIds:     req.VoterIds,
			nbAlts:       req.NbAlts,
			tieBreakRule: req.TieBreakRule,
		}
		var ballot *ballot.Ballot = &b
	case "borda":
		// TODO
	case "approval":
		// TODO
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Unknown rule '%s'", req.Rule)
		w.Write([]byte(msg))
		return
	}

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
	bal := rvsa.getBallotById(req.BallotId)
	if bal == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Pas de ballot avec l'id %v", req.BallotId)
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
		fmt.Fprint(w, "Pas de bollot avec l'id %v", req.BallotId)
	}
	return
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
