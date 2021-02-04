package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type creditRequest struct {
	Address string `json:"address"`
}

type creditResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (f *Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req creditRequest
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "missing address in body", http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		if req.Address = r.FormValue("address"); req.Address == "" {
			http.Error(w, "missing address in query", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sent, err := f.GetTotalSent(req.Address)
	if err != nil {
		sendResponse(w, &creditResponse{
			Status: "failed",
			Error:  "could not get total tokens funded for this account: " + err.Error(),
		})

		return
	}

	if sent >= f.maxCredit {
		log.WithFields(map[string]interface{}{
			"address": req.Address,
			"amount":  fmt.Sprintf("%d%s", f.creditAmount, f.denom),
			"total":   sent + f.creditAmount,
		}).Warnf("tokens not sent: reached maximum credit")

		sendResponse(w, &creditResponse{
			Status: "failed",
			Error:  fmt.Sprintf("account has reached maximum credit allowed per account (%d)", f.maxCredit),
		})

		return
	}

	if err := f.Send(req.Address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"address": req.Address,
		"amount":  fmt.Sprintf("%d%s", f.creditAmount, f.denom),
		"total":   sent + f.creditAmount,
	}).Infof("tokens sent")

	sendResponse(w, &creditResponse{
		Status: "ok",
	})
}

func sendResponse(w http.ResponseWriter, response *creditResponse) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
