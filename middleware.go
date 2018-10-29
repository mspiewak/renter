package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

type errorHandlingMarshalFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func errorHandler(next errorHandlingMarshalFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := next(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if r.Method != http.MethodOptions && (result == nil || (reflect.TypeOf(result).Kind() == reflect.Slice && reflect.ValueOf(result).Len() == 0)) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, []interface{}{})
			return
		}

		response, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(response)
	})
}
