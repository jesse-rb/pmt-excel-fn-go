package pmt

import (
	"context"
	"testing"
	"time"

	"github.com/jesse-rb/pmt-excel-fn-go/pkg/db/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db, err := testutils.SetupTestPostgres(t, ctx)
	defer db.Close()
	require.NoError(t, err)

	repo := NewPostgresRepository(db.Pool)

	timeBeforeCreated := time.Now()
	pmtHistory, err := repo.Create(ctx, 130025034, 0.3293, 24, 42863496)
	timeAfterCreated := time.Now()

	require.NoError(t, err)

	require.GreaterOrEqual(t, pmtHistory.CreatedAt.Compare(timeBeforeCreated), 0)
	require.LessOrEqual(t, pmtHistory.CreatedAt.Compare(timeAfterCreated), 0)

	require.GreaterOrEqual(t, pmtHistory.UpdatedAt.Compare(timeBeforeCreated), 0)
	require.LessOrEqual(t, pmtHistory.UpdatedAt.Compare(timeAfterCreated), 0)

	require.NotEmpty(t, pmtHistory.ID)
	require.Equal(t, int64(130025034), pmtHistory.LoanAmountCents)
	require.Equal(t, float64(0.3293), pmtHistory.InterestRate)
	require.Equal(t, int32(24), pmtHistory.NumPayments)
	require.Equal(t, int64(42863496), pmtHistory.PMTCents)
}
