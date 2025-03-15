package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math"
)


type Inputs struct {
	Efficiency   float64 `json:"efficiency"`
	PowerFactor float64 `json:"powerfactor"`
	Voltage   float64 `json:"voltage"`
	Quantity   float64 `json:"quantity"`
	Ph float64 `json:"ph"`
	Kb      float64 `json:"kb"`
	Tg      float64 `json:"tg"`
}

type Outputs struct {
	Current     float64 `json:"current"`
	Groupusage     float64 `json:"groupusage"`
	Effectiveqty     float64 `json:"effectiveqty"`
	Activepower         float64 `json:"activepower"`
	Reactivepower float64 `json:"reactivepower"`
	Totalpower float64 `json:"totalpower"`
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

	type Inputs struct {
		Efficiency   float64 `json:"efficiency"`
		PowerFactor float64 `json:"powerfactor"`
		Voltage   float64 `json:"voltage"`
		Quantity   float64 `json:"quantity"`
		Ph float64 `json:"ph"`
		Kb      float64 `json:"kb"`
		Tg      float64 `json:"tg"`
	}

	totalpH := req.Quantity * req.Ph
	
	current := (req.Quantity * req.Ph) / (math.Sqrt(3.0) * req.Voltage * req.PowerFactor * req.Efficiency)
	groupUsage := req.Kb * totalpH / totalpH
	effectiveQty := (totalpH * totalpH) / (req.Ph*req.Ph)
	reactivePower := totalpH * req.Kb * req.Tg
	activePower := totalpH * req.Kb
	totalPower := math.Sqrt(activePower * activePower + reactivePower * reactivePower)


	response := Outputs{
		Current: current,
		Groupusage: groupUsage,
		Effectiveqty: effectiveQty,
		Activepower: activePower,
		Reactivepower: reactivePower,
		Totalpower: totalPower,
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
