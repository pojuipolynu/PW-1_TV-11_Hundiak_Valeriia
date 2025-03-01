package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math"
)

type Inputs1 struct {
	Ui   float64 `json:"u"`
	Kz float64 `json:"kz"`
	Time   float64 `json:"time"`
	Sm   float64 `json:"sm"`
	Tm float64 `json:"tm"`
	InputType int `json:"input_type"`
}

type Inputs2 struct {
	Kzu   float64 `json:"kzu"`
	InputType int `json:"input_type"`
}

type Inputs3 struct {
	Rh   float64 `json:"rh"`
	Xh float64 `json:"xh"`
	Rm   float64 `json:"rm"`
	Xm   float64 `json:"xm"`
	InputType int `json:"input_type"`
}

type Outputs1 struct {
	BronK     float64 `json:"bronk"`
	AbbK     float64 `json:"abbk"`
}

type Outputs2 struct {
	Strum3     float64 `json:"strum3"`
}

type Outputs3 struct {
	Ish_3     float64 `json:"ish_3"`
	Ish_3_min     float64 `json:"ish_3_min"`
	Ish_2     float64 `json:"ish_2"`
	Ish_2_min     float64 `json:"ish_2_min"`
	DIsh_3     float64 `json:"dish_3"`
	DIsh_3_min     float64 `json:"dish_3_min"`
	DIsh_2     float64 `json:"dish_2"`
	DIsh_2_min     float64 `json:"dish_2_min"`
}

func calculateCables(req Inputs1) Outputs1 {
	j := 1.0
	if (1000 < req.Tm && req.Tm < 3000) {
		j = 1.6
	} else if (3000 < req.Tm && req.Tm < 5000) {
		j = 1.4
	} else if (5000 < req.Tm) {
		j = 1.2
	}

	bron := (req.Sm / 2) / (math.Sqrt(3.0) * req.Ui) / j
	abb := ((req.Kz * 1000) * math.Sqrt(req.Time) / 92)
	
	response := Outputs1{
		BronK: bron,
		AbbK: abb,
	}
	return response
}

func calculateKz1(req Inputs2) Outputs2 {
	Uc := 10.5
	strum := Uc / (math.Sqrt(3.0) * (Uc*Uc / req.Kzu) + ((Uc / 100) * (Uc*Uc / 6.3)))
	response := Outputs2{
		Strum3: strum,
	}
	return response
}

func calculateKz2(req Inputs3) Outputs3 {
	Xt := (11.1 * 115*115) / (100 * 6.3)

	Xsh := req.Xh + Xt
	Zsh := math.Sqrt(req.Rh*req.Rh + Xsh*Xsh)

	Xsh_min := req.Xm + Xt
	Zsh_min := math.Sqrt(req.Rm*req.Rm + Xsh_min*Xsh_min)

	Ish_3 := (115.0 * math.Pow(10, 3)) / (math.Sqrt(3.0) * Zsh)
	Ish_2 := Ish_3 * (math.Sqrt(3.0) / 2)

	Ish_3_min := (115.0 * math.Pow(10, 3)) / (math.Sqrt(3.0) * Zsh_min)
	Ish_2_min := Ish_3_min * (math.Sqrt(3.0) / 2)

	k := (11.0*11.0) / (115.0*115.0)

	Zsh = math.Sqrt(math.Pow(req.Rh*k, 2)+ math.Pow(Xsh*k, 2))

	Zsh_min = math.Sqrt(math.Pow(req.Rm*k, 2) + math.Pow(Xsh_min*k, 2))

	DIsh_3 := (11.0 * math.Pow(10.0, 3)) / (math.Sqrt(3.0) * Zsh)
	DIsh_2 := Ish_3 * (math.Sqrt(3.0) / 2)

	DIsh_3_min := (11.0 * math.Pow(10, 3)) / (math.Sqrt(3.0) * Zsh_min)
	DIsh_2_min := Ish_3_min * (math.Sqrt(3.0) / 2)
	response := Outputs3{
		Ish_3: Ish_3,     
		Ish_3_min: Ish_3_min,
		Ish_2: Ish_2,    
		Ish_2_min: Ish_2_min, 
		DIsh_3: DIsh_3,     
		DIsh_3_min: DIsh_3_min,
		DIsh_2: DIsh_2,     
		DIsh_2_min: DIsh_2_min,
	}
	return response
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
		response = calculateCables(req)
	case 2:
		var req Inputs2
		jsonBody, _ := json.Marshal(raw)
		json.Unmarshal(jsonBody, &req)
		response = calculateKz1(req)
	case 3:
		var req Inputs3
		jsonBody, _ := json.Marshal(raw)
		json.Unmarshal(jsonBody, &req)
		response = calculateKz2(req)
	default:
		http.Error(w, "Invalid input_type value", http.StatusBadRequest)
		return
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
