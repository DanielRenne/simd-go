package main

import (
	"log"
	"time"

	"github.com/DanielRenne/simd-go/json"
)

type complexObj struct {
	Name       string    `json:"Name"`
	Number     float32   `json:"Number"`
	SomeDate   time.Time `json:"SomeDate"`
	SomeBool   bool      `json:"SomeBool"`
	SomeStruct struct {
		NestedString string    `json:"NestedString"`
		NestedDate   time.Time `json:"NestedDate"`
	} `json:"SomeStruct"`
	SomePointer  *string                  `json:"SomePointer"`
	MapStruct    map[string]interface{}   `json:"MapStruct"`
	IntArray     []int                    `json:"IntArray"`
	Int64Array   []int64                  `json:"Int64Array"`
	Float64Array []float64                `json:"Float64Array"`
	StringArray  []string                 `json:"StringArray"`
	StructArray  []structArrayObj         `json:"StructArray"`
	DateArray    []time.Time              `json:"DateArray"`
	MapArray     []map[string]interface{} `json:"MapArray"`
}

type structArrayObj struct {
	StructString string    `json:"StructString"`
	StructDate   time.Time `json:"StructDate"`
}

func main() {
	obj := complexObj{}
	json.Unmarshal([]byte(`{
	"Name":"John", 
	"Number":1.2, 
	"SomeDate":"2012-04-23T18:25:43.511Z", 
	"SomeBool": true, 
	"SomeStruct":{"NestedString":"HelloNest", "NestedDate":"2013-04-23T19:25:43.511Z"}, 
	"MapStruct":{"MyString":"testMap", "NestedMapStruct":{"MyInt":55555}},
	"SomePointer":"I am a pointer", 
	"IntArray": [1,2,3,4], 
	"Int64Array": [20,30,40],
	"Float64Array": [5.6999,7.6456582,9.1255454],
	"StringArray":["Dan", "Dave", "Bob"],
	"StructArray":[{"StructString":"Hello", "StructDate":"2016-04-23T19:25:43.511Z"},{"StructString":"World", "StructDate":"2018-04-23T19:25:43.511Z"}],
	"DateArray":["2016-04-23T19:25:43.511Z","2018-04-23T19:25:43.511Z"],
	"MapArray":[{"TestObj1":"Value1"}, {"TestObj2":"Value2"}]
	}`), &obj)

	log.Printf("%+v", obj)
}
