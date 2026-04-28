package pmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/money"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Create(ctx context.Context, loanAmountCents int64, interestRate float64, numPayments int32, pmtCents int64) (*PMTHistory, error) {
	args := m.Called(ctx, loanAmountCents, interestRate, numPayments, pmtCents)
	return args.Get(0).(*PMTHistory), args.Error(1)
}

func TestHandlePMT(t *testing.T) {
	type TestCase struct {
		setupMocks           func(*mockStore)
		name                 string
		input                HandlePMTParams
		expectedPMT          string
		expectedResponseCode int
	}

	cases := []TestCase{
		{
			name: "happy - green paths responds with the calculated pmt and http status OK",
			setupMocks: func(store *mockStore) {
				store.
					On(
						"Create",
						mock.Anything, // ctx
						money.ToCents(10_000_000),
						float64(1.2),
						int32(3),
						int64(13_243_781_09),
					).
					Return(&PMTHistory{}, nil)
			},
			input: HandlePMTParams{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  3,
			},
			expectedPMT:          "13243781.09",
			expectedResponseCode: http.StatusOK,
		},
		{
			name: "unhappy - 0 num payments should respond http status Bad Request",
			setupMocks: func(store *mockStore) {
				store.AssertNotCalled(t, "Create")
			},
			input: HandlePMTParams{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  0,
			},
			expectedPMT:          "",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name: "unhappy - db error should respond with http status Internal Server Error",
			setupMocks: func(store *mockStore) {
				store.
					On(
						"Create",
						mock.Anything, // ctx
						money.ToCents(10_000_000),
						float64(1.2),
						int32(3),
						int64(13_243_781_09),
					).
					Return((*PMTHistory)(nil), fmt.Errorf("A db error occured"))
			},
			input: HandlePMTParams{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  3,
			},
			expectedPMT:          "",
			expectedResponseCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Setup mocks
			store := &mockStore{}
			c.setupMocks(store)

			// Setup handler
			handler := NewHTTPHandler(store)
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf(
					"/pmt?loan_amount=%f&interest_rate=%f&num_payments=%d",
					c.input.LoanAmount,
					c.input.InterestRate,
					c.input.NumPayments,
				),
				nil,
			)
			rec := httptest.NewRecorder()

			// Run handler
			handler.HandlePMT(rec, req)
			require.Equal(t, c.expectedResponseCode, rec.Code)

			// Get response body
			var resp HandlePMTResponse
			err := json.NewDecoder(rec.Body).Decode(&resp)
			if c.expectedResponseCode == http.StatusOK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			// Make assertions
			require.Equal(t, c.expectedPMT, resp.PMT)
			store.AssertExpectations(t)
		})
	}
}
