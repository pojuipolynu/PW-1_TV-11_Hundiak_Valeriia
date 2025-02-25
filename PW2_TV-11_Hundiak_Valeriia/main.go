package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math"
)

type Inputs struct {
	CoalWeight   float64 `json:"coalWeight"`
	MazutWeight float64 `json:"mazutWeight"`
	GasWeight   float64 `json:"gasWeight"`
	QCoal   float64 `json:"qCoal"`
}

type Outputs struct {
	KCoal     float64 `json:"kCoal"`
	ECoal     float64 `json:"eCoal"`
	KMazut     float64 `json:"kmazut"`
	EMazut         float64 `json:"emazut"`
	KGas float64 `json:"kGas"`
	EGas float64 `json:"eGas"`
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

	gasDensity := 0.723
	gas := req.GasWeight * gasDensity
	kCoal := math.Pow(10, 6) / req.QCoal * 0.8 * 25.2 / (100 - 1.5) * (1 - 0.985)
	eCoal := math.Pow(10, -6) * kCoal * req.QCoal * req.CoalWeight
	//Розрахунок результатів на основі маси мазуту
	kmazut := math.Pow(10, 6) / 39.48 * 1 * 0.15 / (100 - 0) * (1 - 0.985)
	emazut := math.Pow(10, -6) * kmazut * 39.48 * req.MazutWeight
	//Розрахунок результатів на основі маси газу
	kGas := math.Pow(10, 6) / 33.08 * 0 * 0 / (100 - 0) * (1 - 0.985)
	eGas := math.Pow(10, -6) * kGas * 33.08 * gas


	response := Outputs{
		KCoal: kCoal,
		ECoal: eCoal,
		KMazut: kmazut,
		EMazut: emazut,  
		KGas: kGas,
		EGas: eGas,
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
