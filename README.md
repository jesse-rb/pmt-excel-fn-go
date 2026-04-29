# pmt-excel-fn-go
Calculates monthly repayment amount for a loan.

# Local developemnt

**Requirements**
- docker and docker-compose (docker desktop comes with both docker and docker-compose binaries)
- go v1.26.2
- required go tools `go install tool`
    - tparse: for nice test output summary
    - grpc-gateway binaries and dependencies
- ability to run Makefile

**1. Install tools required by go.mod e.g. grpc-gateway deps**

```
go install tool
```

**2. migratiohs**

- For convenience, migrations are ran automatically when our pmt service starts up, for production this should be changes so migrations are ran on their own
- The local dev postgres DB can be easily reset with `docker-compose down` or using the below make target

```
make docker-down
```

**3. Run services locally (including a postgres db)**

- The below command can be used to run all services in this repo locally (there is only 1 at the moment)
- This uses `docker-compose.yml` to spin up a postgres DB server

```
make docker-build-run
```

**4. Running tests**

running tests uses the tparse go tool, which just changes the output to a nice summary

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
