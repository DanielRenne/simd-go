package main

import "github.com/DanielRenne/simd-go/json"

type complexObj struct {
	Name string `json:"Name"`
}

func main() {
	obj := complexObj{}
	json.Unmarshal([]byte(`{"Name":"John"}`), &obj)
}
