package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EquipmentReliability struct {
	FailureRate        float64 `json:"failure_rate"`
	AverageRepairTime  int     `json:"average_repair_time"`
	Frequency         float64 `json:"frequency"`
	AverageRecoveryTime *int    `json:"average_recovery_time"`
}

var data = map[string]EquipmentReliability{
	"T-110 kV": {0.015, 100, 1.0, intPtr(43)},
	"T-35 kV": {0.02, 80, 1.0, intPtr(28)},
	"T-10 kV (кабельна мережа 10 кV)": {0.005, 60, 0.5, intPtr(10)},
	"T-10 kV (повітряна мережа 10 кV)": {0.05, 60, 0.5, intPtr(10)},
	"B-110 kV (елегазовий)": {0.01, 30, 0.1, intPtr(30)},
	"B-10 kV (малолойний)": {0.02, 15, 0.33, intPtr(15)},
	"B-10 kV (вакуумний)": {0.05, 15, 0.33, intPtr(15)},
	"Збірні шини 10 кV на 1 приєднання": {0.03, 2, 0.33, intPtr(15)},
	"АВ-0,38 кV": {0.05, 20, 1.0, intPtr(15)},
	"ЕД 6,10 кV": {0.1, 50, 0.5, nil},
	"ЕД 0,38 кV": {0.1, 50, 0.5, nil},
	"ПЛ-110 кV": {0.007, 10, 0.167, intPtr(35)},
	"ПЛ-35 кV": {0.02, 8, 0.167, intPtr(35)},
	"ПЛ-10 кV": {0.02, 10, 0.167, intPtr(35)},
	"КЛ-10 кV (траншея)": {0.03, 44, 1.0, intPtr(9)},
	"КЛ-10 кV (кабельний канал)": {0.005, 18, 1.0, intPtr(9)},
}

type Inputs1 struct {
	Omega     float64 `json:"omega"`
	Tsvalue   float64 `json:"tsvalue"`
	Pmvalue   float64 `json:"pmvalue"`
	Tmvalue   float64 `json:"tmvalue"`
	Zavarvalue float64 `json:"zavarvalue"`
	Zplanvalue float64 `json:"zplanvalue"`
	Kpvalue   float64 `json:"kpvalue"`
	InputType int     `json:"input_type"`
}

type Inputs2 struct {
	Amounts map[string]int `json:"amounts"`
}

type Outputs1 struct {
	MWnedA float64 `json:"mwneda"`
	MWnedP float64 `json:"mwnedp"`
	MZ     float64 `json:"mz"`
}

type Outputs2 struct {
	Woc  float64 `json:"woc"`
	Tvoc float64 `json:"tvoc"`
	Kaoc float64 `json:"kaoc"`
	Kpoc float64 `json:"kpoc"`
	Wdk  float64 `json:"wdk"`
	Wds  float64 `json:"wds"`
}

func calculateDeficit(req Inputs1) Outputs1 {
	mWnedAValue := req.Omega * req.Tsvalue * req.Pmvalue * req.Tmvalue
	mWnedPValue := req.Kpvalue * req.Pmvalue * req.Tmvalue
	mZValue := req.Zavarvalue * mWnedAValue + req.Zplanvalue * mWnedPValue

	return Outputs1{
		MWnedA: mWnedAValue,
		MWnedP: mWnedPValue,
		MZ:     mZValue,
	}
}

func calculateReliability(req Inputs2) Outputs2 {
	woc := 0.0
	tvoc := 0.0
	for key, amount := range req.Amounts {
		if entry, exists := data[key]; exists && amount > 0 {
			woc += float64(amount) * entry.FailureRate
			tvoc += float64(amount) * float64(entry.AverageRepairTime) * entry.FailureRate
		}
	}

	if woc > 0 {
		tvoc /= woc
	}
	kaoc := (tvoc * woc) / 8760
	kpoc := 1.2 * 43 / 8760
	wdk := 2 * woc * (kaoc + kpoc)
	wds := wdk + 0.02

	return Outputs2{
		Woc:  woc,
		Tvoc: tvoc,
		Kaoc: kaoc,
		Kpoc: kpoc,
		Wdk:  wdk,
		Wds:  wds,
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	inputTypeFloat, ok := raw["input_type"].(float64)
	if !ok {
		http.Error(w, "Invalid or missing input_type", http.StatusBadRequest)
		return
	}
	inputType := int(inputTypeFloat)

	var response interface{}
	switch inputType {
	case 1:
		var req Inputs1
		jsonBody, _ := json.Marshal(raw)
		json.Unmarshal(jsonBody, &req)
		response = calculateDeficit(req)
	case 2:
		var req Inputs2
		jsonBody, _ := json.Marshal(raw)
		json.Unmarshal(jsonBody, &req)
		response = calculateReliability(req)
	default:
		http.Error(w, "Invalid input_type value", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func intPtr(i int) *int {
	return &i
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	http.Handle("/", http.FileServer(http.Dir("./pages")))
	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
