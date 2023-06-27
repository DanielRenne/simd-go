package json

import (
	"encoding/json"
	"log"

	"github.com/minio/simdjson-go"
)

// func Marshall(v any) (data []byte, err error) {
// 	json.Marshal(v any)
// }

func Unmarshal(data []byte, v any) (err error) {

	if simdjson.SupportedCPU() {
		pj, errParse := simdjson.Parse(data, nil)
		if errParse != nil {
			err = errParse
			return

		}
		log.Printf("%+v", pj)
		return
	}

	err = json.Unmarshal(data, v)
	return
}
