package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

// @ TYPES
// recieves fields attach any values to it
type Envelop map[string] interface{} //* used mostly with json responses

// ! Functions that handle json reponses {data --> data that would be passed to process response}
func ReadJson(w http.ResponseWriter,r *http.Request,data interface{}) error {
	
	// # configuring how much data could be processed
	maxByteIngress := 1048576
	r.Body = http.MaxBytesReader(w,r.Body,int64(maxByteIngress))

	// # Dealin with incoming data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	err = decoder.Decode(&struct{}{})
	if err != nil {
		return errors.New("Body must have only a single json object")
	}
	return nil
}

func WriteJson(w http.ResponseWriter,status int,data interface{},headers ...http.Header) error {
	output,err := json.MarshalIndent(data,"", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(status)
	_,err = w.Write(output)
	if err != nil {
		return err
	}
	return  nil	
}

// func ErrorJson(w http.ResponseWriter,err error,status ...int) {
// 	statusCode := http.StatusBadRequest
// 	if len(status) > 0 {
// 		statusCode = status[0]
// 	}

// }