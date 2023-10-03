package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sprint/internal/storage"
	"sprint/internal/utils"

	"go.uber.org/zap"
)

type DataReq struct {
	URL string `json:"url"`
}

type DataResp struct {
	Result string `json:"result"`
}

func HandlerPostAPI(w http.ResponseWriter, r *http.Request, baseAddress, file string, log *zap.Logger) {
	if r.URL.Path != "/api/shorten" || utils.ValidContentType(r.Header.Get("Content-Type"), "application/json") != nil {
		log.Info("PostAPI not Path or not Content-Type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	var data DataReq
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Info("PostAPI not body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		log.Info("PostAPI not byte to json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if data.URL == "" {
		log.Info("PostAPI not url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.LongToShort(data.URL, file)
	dataResp := DataResp{Result: baseAddress + "/" + storage.GetDB(string(data.URL))}
	resp, err := json.Marshal(dataResp)
	if err != nil {
		log.Info("PostAPI not json to byte")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
