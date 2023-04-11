package utils

import (
	"encoding/json"
	"net/http"
)

func RenderJson(r *http.Request, w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8 ")
	var body []byte
	if res != nil {
		var err error
		body, err = json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(statusCode)
	if len(body) > 0 {
		w.Write(body)
	}
}
