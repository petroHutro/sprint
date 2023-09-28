package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sprint/internal/storage"
	"sprint/internal/utils"
)

type DataReq struct {
	URL string `json:"url"`
}

type DataResp struct {
	Result string `json:"result"`
}

func HandlerPostAPI(w http.ResponseWriter, r *http.Request, baseAddress string) {
	if r.URL.Path != "/api/shorten" || utils.ValidContentType(r.Header.Get("Content-Type"), "application/json") != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	var data DataReq
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fmt.Println(data)
	if data.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.LongToShort(data.URL)
	var dataResp DataResp = DataResp{Result: baseAddress + "/" + storage.GetDB(string(data.URL))}
	resp, err := json.Marshal(dataResp)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
