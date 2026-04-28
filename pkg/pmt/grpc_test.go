package pmt

import (
	"context"
	"fmt"
	"testing"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/money"
	"github.com/jesse-rb/pmt-excel-fn-go/protogen"
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
	ctx := context.Background()

	type TestCase struct {
		setupMocks    func(*mockStore)
		name          string
		input         protogen.PMTRequest
		expectedPMT   string
		expectedError bool
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
			input: protogen.PMTRequest{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  3,
			},
			expectedPMT:   "13243781.09",
			expectedError: false,
		},
		{
			name: "unhappy - 0 num payments should respond http status Bad Request",
			setupMocks: func(store *mockStore) {
				store.AssertNotCalled(t, "Create")
			},
			input: protogen.PMTRequest{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  0,
			},
			expectedPMT:   "",
			expectedError: true,
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
			input: protogen.PMTRequest{
				LoanAmount:   10_000_000,
				InterestRate: 1.2,
				NumPayments:  3,
			},
			expectedPMT:   "",
			expectedError: true,
		},
	}

	for i := range cases {
		c := &cases[i]
		t.Run(c.name, func(t *testing.T) {
			// Setup mocks
			store := &mockStore{}
			c.setupMocks(store)

			// Setup handler
			handler := &GRPCHandler{
				store: store,
			}

			// Run handler
			resp, err := handler.CalculatePMT(ctx, &c.input)

			if !c.expectedError {
				require.NoError(t, err)
				require.Equal(t, c.expectedPMT, resp.Pmt)
			} else {
				require.Error(t, err)
			}

			// Make assertions
			store.AssertExpectations(t)
		})
	}
}
