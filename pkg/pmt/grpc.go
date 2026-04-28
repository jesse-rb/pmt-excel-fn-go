package pmt

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/money"
	"github.com/jesse-rb/pmt-excel-fn-go/protogen"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	protogen.UnimplementedPMTServiceServer
	service *Service
	store   Store
}

func NewGRPCHandler(grpc *grpc.Server, service *Service, store Store) {
	grpcHandler := &GRPCHandler{
		service: service,
		store:   store,
	}

	protogen.RegisterPMTServiceServer(grpc, grpcHandler)
}

func (h *GRPCHandler) CalculatePMT(ctx context.Context, req *protogen.PMTRequest) (*protogen.PMTResponse, error) {
	// Check input
	if req.NumPayments == 0 {
		return nil, ErrZeroNumPayments
	}

	// Calc PMT
	loanAmountCents := money.ToCents(req.LoanAmount)
	pmtCents, err := h.service.CalcPMT(loanAmountCents, req.InterestRate, req.NumPayments)

	// Store record
	_, err = h.store.Create(ctx, loanAmountCents, req.InterestRate, req.NumPayments, pmtCents)
	if err != nil {
		slog.Error("failed to store pmt history", "err", err)
		return nil, ErrInternal
	}

	// Respond
	pmt := money.FromCents(pmtCents)
	resp := protogen.PMTResponse{
		Pmt: fmt.Sprintf("%.2f", pmt),
	}

	return &resp, err
}
