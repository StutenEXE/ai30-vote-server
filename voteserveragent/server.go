package voteserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"https://github.com/StutenEXE/ai30-vote-server/td5"
)

type RestVoteServerAgent struct {
	sync.Mutex
	id   string
	port string
}

func NewRestVoteServerAgent(port string) *RestVoteServerAgent {
	return &RestVoteServerAgent{id: port, port: port}
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

fun(*RestVoteServerAgent) decodeVoteRequest(r *http.Request) (req td5.VoteRequest, err error){
	
}

// func (rsa *RestServerAgent) doCalc(w http.ResponseWriter, r *http.Request) {
// 	// mise à jour du nombre de requêtes
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	rsa.reqCount++

// 	// vérification de la méthode de la requête
// 	if !rsa.checkMethod("POST", w, r) {
// 		return
// 	}

// 	// décodage de la requête
// 	req, err := rsa.decodeRequest(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	// traitement de la requête
// 	var resp rad.Response

// 	switch req.Operator {
// 	case "*":
// 		resp.Result = req.Args[0] * req.Args[1]
// 	case "+":
// 		resp.Result = req.Args[0] + req.Args[1]
// 	case "-":
// 		resp.Result = req.Args[0] - req.Args[1]
// 	default:
// 		w.WriteHeader(http.StatusNotImplemented)
// 		msg := fmt.Sprintf("Unkonwn command '%s'", req.Operator)
// 		w.Write([]byte(msg))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	serial, _ := json.Marshal(resp)
// 	w.Write(serial)
// }

// func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
// 	if !rsa.checkMethod("GET", w, r) {
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	rsa.Lock()
// 	defer rsa.Unlock()
// 	serial, _ := json.Marshal(rsa.reqCount)
// 	w.Write(serial)
// }

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

	req

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
	req
}

func (rvsa *RestVoteServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rvsa.addBallot)
	mux.HandleFunc("/vote", rvsa.addVote)
	mux.HandleFunc("/result", rvsa.giveResult)

	// création du serveur http
	s := &http.Server{
		Addr:           rvsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
