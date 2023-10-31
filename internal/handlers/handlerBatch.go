package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strconv"
)

type DataReqBatch struct {
	ID   string `json:"correlation_id"`
	Long string `json:"original_url"`
}

type DataRespBatch struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

func HandlerPostBatch(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	var buf bytes.Buffer
	var data []DataReqBatch

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("PostBatch not body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logger.Error("PostBatch not byte to json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var longs []string

	for _, item := range data {
		if item.ID == "" || item.Long == "" {
			logger.Error("PostBatch not correlation_id or original_url: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		longs = append(longs, item.Long)
	}

	statusCode := http.StatusCreated
	id, err := strconv.Atoi(r.Header.Get("User_id"))
	if err != nil {
		logger.Error("bad user id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := db.SetAll(r.Context(), longs, id); err != nil {
		var repErr *storage.RepError
		if errors.As(err, &repErr) && repErr.Repetition {
			statusCode = http.StatusConflict
			logger.Error("long already db :%v", err)
		} else {
			logger.Error("cannot set all: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var dataResp []DataRespBatch

	for _, item := range data {
		short, _ := db.GetShort(r.Context(), item.Long) //!!!!!!!!!!!!!!!!!!!!
		dataResp = append(dataResp, DataRespBatch{
			ID:    item.ID,
			Short: baseAddress + "/" + short})
	}

	resp, err := json.Marshal(dataResp)
	if err != nil {
		logger.Error("PostAPI not json to byte :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
