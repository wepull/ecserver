package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	PORT = ":8081"
)

var (
	candidates = []*Candidate{
		{
			Name:     "Roost",
			ImageUrl: "https://roost.ai/hubfs/Roost-V3.svg",
		},
		{
			Name:     "Docker",
			ImageUrl: "https://logos-world.net/wp-content/uploads/2021/02/Docker-Logo.png",
		},
 		{
 			Name:     "MiniKube",
 			ImageUrl: "https://miro.medium.com/max/400/0*KzqL3xqmXzV5PPjX.png",
 		},
// 		{
// 		 	Name:     "K3D",
// 			ImageUrl: "https://www.suse.com/c/wp-content/uploads/2021/02/K3D-Blog-Graphic-20210209.png",
// 		},
		{
			Name: "Rancher",
			ImageUrl: "https://rancher.com/assets/img/logos/rancher-logo-horiz-color.svg",
		},
	}
)

// Basic payload to be sent in response to API request
type BasicResponse struct {
	Status  int
	Message string
}

// Payload to send back all candidates
type AllCandidatesResponse struct {
	Candidates []*Candidate
}

type Candidate struct {
	Name     string
	ImageUrl string
}

func writeBasicResponse(w http.ResponseWriter, resp *BasicResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Println("error marshaling response to vote request. error: ", err)
	}
	w.Write(respJson)
}

func writeAllCandidatesResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respJson, err := json.Marshal(&AllCandidatesResponse{
		Candidates: candidates,
	})
	if err != nil {
		log.Println("error marshaling response to vote request. error: ", err)
	}
	w.Write(respJson)
}

func getAllCandidates(w http.ResponseWriter, r *http.Request) {
	writeAllCandidatesResponse(w)
}

func addCandidate(w http.ResponseWriter, r *http.Request) {
	newCandidate := &Candidate{}
	err := json.NewDecoder(r.Body).Decode(newCandidate)
	if err != nil {
		resp := &BasicResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		writeBasicResponse(w, resp)
		return
	}

	isCandidatePresent := false
	for i, ca := range candidates {
		if newCandidate.Name == ca.Name {
			isCandidatePresent = true
			candidates[i].Name = newCandidate.Name
			candidates[i].ImageUrl = newCandidate.ImageUrl
		}
	}

	if !isCandidatePresent {
		candidates = append(candidates, newCandidate)
	}

	writeAllCandidatesResponse(w)
}

func deleteCandidate(w http.ResponseWriter, r *http.Request) {
	reqDeleteCandidate := &Candidate{}
	err := json.NewDecoder(r.Body).Decode(reqDeleteCandidate)
	if err != nil {
		resp := &BasicResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		writeBasicResponse(w, resp)
		return
	}

	isCandidatePresent := false
	updatedCandidates := []*Candidate{}
	for i, ca := range candidates {
		if reqDeleteCandidate.Name != ca.Name {
			updatedCandidates = append(updatedCandidates, candidates[i])
		} else {
			isCandidatePresent = true
		}
	}
	candidates = updatedCandidates

	if !isCandidatePresent {
		resp := &BasicResponse{
			Status:  http.StatusBadRequest,
			Message: "Candidate Not Found",
		}
		writeBasicResponse(w, resp)
		return
	}

	writeAllCandidatesResponse(w)
}

func handleInvalidMethod(w http.ResponseWriter, r *http.Request) {
	resp := &BasicResponse{
		Status:  http.StatusMethodNotAllowed,
		Message: "No such endpoint",
	}
	writeBasicResponse(w, resp)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case http.MethodGet:
		getAllCandidates(w, r)

	case http.MethodPost:
		addCandidate(w, r)

	case http.MethodDelete:
		deleteCandidate(w, r)

	case http.MethodOptions:
		return

	default:
		handleInvalidMethod(w, r)
	}
}

func main() {
	http.HandleFunc("/", serveRoot)
	log.Printf("Starting a server at port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
