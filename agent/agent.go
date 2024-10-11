package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"td5"
	"time"
)

var url string = "http://localhost:8080"

type Agent struct {
	id string
}

func NewAgent(id string) *Agent {
	return &Agent{id: id}
}

// J'ai mis un interface pour tes les types
// de requete afin de faire tous les post possible
func (vca *Agent) makeRequest(endpoint string, method string, mesdata interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", url, endpoint)

	data, err := json.Marshal(mesdata)
	if err != nil {
		return nil, err
	}

	// Création de la requête HTTP
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Définition des headers
	req.Header.Set("Content-Type", "application/json")

	// Envoi de la requête
	// JE sais pas trop trop quoi mettre en timeout
	/// TODOO
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	return res, err
}

func (vca *Agent) Ballot(req td5.BallotRequest) (string, string, error) {
	// Envoi de la requête de création de bulletin
	resp, err := vca.makeRequest("new_ballot", "POST", req)
	if err != nil {
		return resp.Status, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return resp.Status, "", fmt.Errorf("failed to create ballot: %s", resp.Status)
	}

	var id string
	err = json.NewDecoder(resp.Body).Decode(&id)
	if err != nil {
		return resp.Status, "", err
	}
	return resp.Status, id, nil
}

func (vca *Agent) Vote(req td5.VoteRequest) (string, error) {
	// Envoi de la requête de vote
	resp, err := vca.makeRequest("vote", "POST", req)
	if err != nil {
		return resp.Status, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.Status, fmt.Errorf("failed to vote: %s", resp.Status)
	}

	return resp.Status, nil
}

func (vca *Agent) Result(req td5.ResultRequest) (string, td5.ResultResponse, error) {
	var result td5.ResultResponse

	// Envoi de la requête pour obtenir le résultat
	resp, err := vca.makeRequest("result", "POST", req)
	if err != nil {
		return resp.Status, result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.Status, result, fmt.Errorf("failed to get result: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return resp.Status, result, err
	}

	return resp.Status, result, nil
}
