package json

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

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

		// Iterate each top level element.
		err = pj.ForEach(func(i simdjson.Iter) error {
			fmt.Println("Got iterator for type:", i.Type())

			return unmarshalRecursive(i, reflect.ValueOf(v).Elem())

			// element, err := i.FindElement(nil, "Image", "URL")
			// if err == nil {
			// 	value, _ := element.Iter.StringCvt()
			// 	fmt.Println("Found element:", element.Name, "Type:", element.Type, "Value:", value)
			// }
			// return nil
		})

		// log.Printf("%+v", pj)
		return
	}

	err = json.Unmarshal(data, v)
	return
}

func unmarshalRecursive(i simdjson.Iter, structValue reflect.Value) error {

	structType := structValue.Type()

	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		fieldValue := structValue.Field(j)

		tagName := field.Tag.Get("json")
		// fmt.Println(tagName)
		// fmt.Println(fieldValue.Kind())

		element, err := i.FindElement(nil, tagName)
		if err == nil {
			switch fieldValue.Kind() {
			case reflect.Ptr:
				log.Println("Need to Implement Ptr")
			// 	if fieldValue.IsNil() {
			// 		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			// 	}

			// 	unmarshalRecursive(element.Iter, fieldValue)

			case reflect.String:
				val, err := element.Iter.StringCvt()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(val))
				}
			case reflect.Int:
				val, err := element.Iter.Int()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(int(val)))
				}
			case reflect.Int64:
				val, err := element.Iter.Int()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(val))
				}
			case reflect.Float32:
				val, err := element.Iter.Float()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(float32(val)))
				}
			case reflect.Float64:
				val, err := element.Iter.Float()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(val))
				}
			case reflect.Bool:
				val, err := element.Iter.Bool()
				if err == nil {
					fieldValue.Set(reflect.ValueOf(val))
				}
			case reflect.Map:
				target := make(map[string]interface{})
				obj, err := element.Iter.Object(nil)
				if err == nil {
					obj.ForEach(func(key []byte, i simdjson.Iter) {
						inf, err := i.Interface()
						if err == nil {
							target[string(key)] = inf
						}

					}, nil)
					fieldValue.Set(reflect.ValueOf(target))
				}
			case reflect.Struct:
				if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					val, err := element.Iter.StringCvt()
					if err == nil {
						timeValue, err := time.Parse(time.RFC3339, val)
						if err == nil {
							fieldValue.Set(reflect.ValueOf(timeValue))
						}
					}
				} else {
					unmarshalRecursive(element.Iter, fieldValue)
				}
			case reflect.Slice:

				underLyingType := reflect.TypeOf(fieldValue.Interface()).Elem()

				t := reflect.SliceOf(fieldValue.Type())
				if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf([]int{}) {
					valArray, err := element.Iter.Array(nil)
					if err == nil {
						val, err := valArray.AsInteger()
						if err == nil {

							intSlice := make([]int, len(val))
							for i, v := range val {
								intSlice[i] = int(v)
							}
							fieldValue.Set(reflect.ValueOf(intSlice))
						}
					} else {
						log.Println("Failed to parse Array of int:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf([]int64{}) {
					valArray, err := element.Iter.Array(nil)
					if err == nil {
						val, err := valArray.AsInteger()
						if err == nil {
							fieldValue.Set(reflect.ValueOf(val))
						}
					} else {
						log.Println("Failed to parse Array of int64:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf([]float64{}) {
					valArray, err := element.Iter.Array(nil)
					if err == nil {
						val, err := valArray.AsFloat()
						if err == nil {
							fieldValue.Set(reflect.ValueOf(val))
						}
					} else {
						log.Println("Failed to parse Array of int64:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf([]string{}) {
					valArray, err := element.Iter.Array(nil)
					if err == nil {
						val, err := valArray.AsStringCvt()
						if err == nil {
							fieldValue.Set(reflect.ValueOf(val))
						}
					} else {
						log.Println("Failed to parse Array of string:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf([]time.Time{}) {
					valArray, err := element.Iter.Array(nil)
					if err == nil {
						val, err := valArray.AsStringCvt()
						if err == nil {

							timeSlice := make([]time.Time, len(val))
							for i, v := range val {
								timeValue, err := time.Parse(time.RFC3339, v)
								if err == nil {
									timeSlice[i] = timeValue
								}
							}
							fieldValue.Set(reflect.ValueOf(timeSlice))
						}
					} else {
						log.Println("Failed to parse Array of time strings:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && underLyingType.Kind() == reflect.Map {

					valArray, err := element.Iter.Array(nil)
					if err == nil {

						// // Create a new slice with the type of the array property
						sliceType := reflect.SliceOf(underLyingType)
						newSlice := reflect.MakeSlice(sliceType, 0, 0)

						valArray.ForEach(func(arrayIter simdjson.Iter) {

							target := make(map[string]interface{})
							obj, err := arrayIter.Object(nil)
							if err == nil {
								obj.ForEach(func(key []byte, i simdjson.Iter) {
									inf, err := i.Interface()
									if err == nil {
										target[string(key)] = inf
									}

								}, nil)
							}

							newSlice = reflect.Append(newSlice, reflect.ValueOf(target))
						})

						fieldValue.Set(newSlice)

					} else {
						log.Println("Failed to parse Array of struct:  " + err.Error())
					}
				} else if t.Kind() == reflect.Slice && underLyingType.Kind() == reflect.Struct {

					valArray, err := element.Iter.Array(nil)
					if err == nil {

						// // Create a new slice with the type of the array property
						sliceType := reflect.SliceOf(underLyingType)
						newSlice := reflect.MakeSlice(sliceType, 0, 0)

						valArray.ForEach(func(arrayIter simdjson.Iter) {
							obj := reflect.New(underLyingType).Elem()
							unmarshalRecursive(arrayIter, obj)
							newSlice = reflect.Append(newSlice, obj)
						})

						fieldValue.Set(newSlice)

					} else {
						log.Println("Failed to parse Array of struct:  " + err.Error())
					}
				}
			}

		} else {
			log.Println("Failed to find element:  " + err.Error())
		}
	}

	return nil
}

// func unmarshalRecursive(node *simdjson.Object, structValue reflect.Value) error {
// 	structType := structValue.Type()

// 	for i := 0; i < structType.NumField(); i++ {
// 		field := structType.Field(i)
// 		fieldValue := structValue.Field(i)

// 		tagName := field.Tag.Get("json")
// 		jsonValue := node.Fin Get(tagName)

// 		if jsonValue == nil {
// 			continue
// 		}

// 		switch fieldValue.Kind() {
// 		case reflect.Ptr:
// 			if fieldValue.IsNil() {
// 				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
// 			}
// 			ptrValue := fieldValue.Elem()
// 			if err := unmarshalRecursive(jsonValue.GetObject(), ptrValue); err != nil {
// 				return err
// 			}

// 		case reflect.Struct:
// 			if err := unmarshalRecursive(jsonValue.GetObject(), fieldValue); err != nil {
// 				return err
// 			}

// 		default:
// 			if err := jsonValue.Unmarshal(&fieldValue); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
