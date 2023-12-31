package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

type Data []string

func HandlerDelete(w http.ResponseWriter, r *http.Request, delChan chan storage.QueryDelete) {
	var data Data

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("PostAPI not body :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logger.Error("PostAPI not byte to json :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.Header.Get("User_id")

	if id == "" {
		logger.Error("bad user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, short := range data {
		delChan <- storage.QueryDelete{ID: id, Data: short}
	}

	w.WriteHeader(http.StatusAccepted)
}
