package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Inputs struct {
	Carbon   float64 `json:"carbon"`
	Hydrogen float64 `json:"hydrogen"`
	Sulfur   float64 `json:"sulfur"`
	Oxygen   float64 `json:"oxygen"`
	Moisture float64 `json:"moisture"`
	Ash      float64 `json:"ash"`
	Venadii      float64 `json:"venadii"`
	Q      float64 `json:"q"`
}

type Outputs struct {
	CoefficientWtoD     float64 `json:"coefficientWtoD"`
	CoefficientWtoC     float64 `json:"coefficientWtoC"`
	HeatWorkingMass     float64 `json:"heatWorkingMass"`
	HeatDryMass         float64 `json:"heatDryMass"`
	HeatCombustibleMass float64 `json:"heatCombustibleMass"`
	WorkingMassCarbon float64 `json:"workingMassCarbon"`
	WorkingMassSulfur float64 `json:"workingMassSulfur"`
	WorkingMassOxygen float64 `json:"workingMassOxygen"`
	WorkingMassVenadii float64 `json:"workingMassVenadii"`
	WorkingMassAsh float64 `json:"workingMassAsh"`
	WorkingMassHydrogen float64 `json:"workingMassHydrogen"`
	QWorkingMass     	float64 `json:"qWorkingMass"`
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

	coefficientWtoD := 100 / (100 - req.Moisture)
	coefficientWtoC := 100 / (100 - req.Moisture - req.Ash)

	heatWorkingMass := (339*req.Carbon + 1030*req.Hydrogen - 108.8*(req.Oxygen-req.Sulfur) - 25*req.Moisture) / 1000
	heatDryMass := (heatWorkingMass + 0.025*req.Moisture) * 100 / (100 - req.Moisture)
	heatCombustibleMass := (heatWorkingMass + 0.025*req.Moisture) * 100 / (100 - req.Moisture - req.Ash)

	m := (100 - req.Moisture - req.Ash) / 100

	workingMassCarbon := req.Carbon * m
	workingMassSulfur := req.Sulfur * m
	workingMassOxygen := req.Oxygen * m
	workingMassVenadii := req.Venadii * m
	workingMassAsh := req.Ash * m
	workingMassHydrogen := req.Hydrogen * m

	qWorkingMass := req.Q * ((100 - req.Moisture - req.Ash) / 100) - 0.025 * req.Moisture


	response := Outputs{
		CoefficientWtoD:     coefficientWtoD,
		CoefficientWtoC:     coefficientWtoC,
		HeatWorkingMass:     heatWorkingMass,
		HeatDryMass:         heatDryMass,
		HeatCombustibleMass: heatCombustibleMass,
		WorkingMassCarbon : workingMassCarbon,
		WorkingMassSulfur : workingMassSulfur,
		WorkingMassOxygen : workingMassOxygen,
		WorkingMassVenadii : workingMassVenadii,
		WorkingMassAsh : workingMassAsh,
		WorkingMassHydrogen : workingMassHydrogen,
		QWorkingMass: qWorkingMass,
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
