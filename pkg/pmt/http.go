package pmt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jesse-rb/pmt-excel-fn-go/pkg/money"
)

type HTTPHandler struct {
	store Store
}

func NewHTTPHandler(store Store) *HTTPHandler {
	return &HTTPHandler{
		store: store,
	}
}

type HandlePMTParams struct {
	LoanAmount   float64 `schema:"loan_amount"`
	InterestRate float64 `schema:"Interest_rate"`
	NumPayments  int32   `schema:"num_payments"`
}

type HandlePMTResponse struct {
	PMT string `json:"pmt"`
}

var decoder = schema.NewDecoder()

func (h HTTPHandler) HandlePMT(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// Scan query params
	var params HandlePMTParams
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		slog.Error("invalid request: failed to decode query params", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if params.NumPayments == 0 {
		slog.Error("invalid request: 0 num payments")
		http.Error(w, "num_payments must be not be 0", http.StatusBadRequest)
		return
	}

	// Calc pmt
	loanAmountCents := money.ToCents(params.LoanAmount)
	pmtCents, err := CalcPMT(loanAmountCents, params.InterestRate, params.NumPayments)
	if err != nil {
		slog.Error("failed to calc PMT", "err", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Store record
	_, err = h.store.Create(context.Background(), loanAmountCents, params.InterestRate, params.NumPayments, pmtCents)
	if err != nil {
		slog.Error("failed to store pmt history", "err", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Respond
	resp := &HandlePMTResponse{
		PMT: fmt.Sprintf("%.2f", money.FromCents(pmtCents)),
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json response", "err", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}
