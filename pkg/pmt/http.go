package pmt

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/schema"
	// "github.com/jesse-rb/pmt-excel-fn-go/pkg/money"
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

var decoder = schema.NewDecoder()

func (h HTTPHandler) HandlePMT(w http.ResponseWriter, r *http.Request) {
	// Scan query params
	var params HandlePMTParams
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		slog.Error("failed to decode query params", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Calc pmt
	// loanAmountCents := money.ToCents(params.LoanAmount)
	// pmt := CalcPMT(loanAmountCents, params.InterestRate, params.NumPayments)
}
