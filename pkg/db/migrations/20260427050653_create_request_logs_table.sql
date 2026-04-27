-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pmt_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_amount_cents BIGINT NOT NULL,
    interest_rate DECIMAL(10,6) NOT NULL, -- 10 digits significant digits, 6 fractional digits after decimal point
    num_payments INTEGER NOT NULL,
    pmt_cents BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pmt_histories;
-- +goose StatementEnd
