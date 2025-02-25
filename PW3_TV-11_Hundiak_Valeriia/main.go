package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math"
)

type Inputs struct {
	P   float64 `json:"p"`
	Q1 float64 `json:"q1"`
	Q2   float64 `json:"q2"`
	B   float64 `json:"b"`
}

type Outputs struct {
	Q1Result_1     float64 `json:"q1result_1"`
	Q1Result_2     float64 `json:"q1result_2"`
	Q1Result_3     float64 `json:"q1result_3"`
	Q2Result_1     float64 `json:"q2result_1"`
	Q2Result_2	   float64 `json:"q2result_2"`
	Q2Result_3	   float64 `json:"q2result_3"`
}

func calculateEnergyShare(P, q float64) float64 {
	delta := P * 0.05
	lowerBound := P-delta
	upperBound :=P+delta
	step := 0.001
	integral := 0.0
	for i:=lowerBound; i < upperBound; i+= step {
	  pd := (1/(q*math.Sqrt(2*math.Pi)))*math.Exp(-(math.Pow(i-P, 2)/2*math.Pow(q, 2)))
	  integral += pd* step
	}
	return integral
}

func calculateProfitAndPenalty(P, cost, q float64) (profit, penalty, diff float64) {
	energyShare := calculateEnergyShare(P, q)
	energyWithoutImbalance := P * 24 * energyShare
	profit = energyWithoutImbalance * cost * 1000
	energyWithImbalance := P * 24 * (1 - energyShare)
	penalty = energyWithImbalance * cost * 1000
	diff = profit-penalty 
	return
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error occurred", http.StatusMethodNotAllowed)
		return
	}

	var req Inputs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error occurred", http.StatusBadRequest)
		return
	}

	q1result_1, q1result_2, q1result_3 := calculateProfitAndPenalty(req.P, req.B, req.Q1)
	q2result_1, q2result_2, q2result_3 := calculateProfitAndPenalty(req.P, req.B, req.Q2)


	response := Outputs{
		Q1Result_1: q1result_1,
		Q1Result_2: q1result_2,
		Q1Result_3: q1result_3,
		Q2Result_1: q2result_1,
		Q2Result_2: q2result_2,
		Q2Result_3: q2result_3,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	http.Handle("/", http.FileServer(http.Dir("./pages")))
	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
