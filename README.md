# pmt-excel-fn-go
Calculates monthly repayment amount for a loan.

# Local developemnt

**Requirements**
- docker and docker-compose (docker desktop comes with both docker and docker-compose binaries)
- go v1.26.2
- tparse `go install github.com/mfridman/tparse@v0.18` for nice test output summary
- ability to run Makefile

# Running and testing services locally

**Install tools required by go.mod e.g. grpc-gateway deps**

```
go install tool
```

**migratiohs**
- For convenience in local dev, migrations are ran automatically when our pmt service starts up

**Run services locally (including a postgres db)**

- The below command can be used to run all services in this repo locally (there is only 1 at the moment)
- This uses `docker-compose.yml` to spin up a postgres DB server

```
make docker-build-run
```

**resetting local dev DB**

- The local dev postgres DB can be easily reset with `docker-compose down` or using the below make target

```
make docker-down
```

**Running tests**

```
make test [pkg?] [run?]

e.g.

make test
make test PKG='./pkg/money ./pkg/pmt'
make test PKG='./pkg/money' RUN='TestToCents'
```

# Services

## PMT (excel function)

**usage examples**

http (calls proxy grpc server)
```
curl -X POST http://localhost:8080/pmt \
  -H "Content-Type: application/json" \
  -d '{
    "loan_amount": 100000,
    "interest_rate": 0.05,
    "num_payments": 360
  }'
```
