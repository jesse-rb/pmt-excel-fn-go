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

**using curl -- http (calls proxy grpc server)**
```
curl -X POST http://localhost:8080/pmt \
  -H "Content-Type: application/json" \
  -d '{
    "loan_amount": 100000,
    "interest_rate": 0.05,
    "num_payments": 3
  }'
```

**using grpcurl -- grpc**

`go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

```
grpcurl -plaintext -d '{
    "loan_amount": 100000,
    "interest_rate": 0.05,
    "num_payments": 3
  }' \
  -import-path ./protobuf \
  -proto pmt.proto \
  localhost:9090 PMTService/CalculatePMT
```

# k8s

**Requirements**

- A k8s cluster e.g. locally via minikube
- kubectl k8s cli tool

**Start k8s cluster (assuming minikube)**

Start a local k8s cluster via minikube and docker
```
minikube start --driver=docker
```

When finished, you can stop or destroy the minikube cluster
```
minikube stop
minikube destroy
```

Check minikube is ready:

```
> kubectl get nodes                                                                                                                                            ?  11s ?? minikube ?  16:16:46
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   38s   v1.35.1
```

## k8s Makefile targets

Below are "## Detailed k8s deployment instructions" however, the following Makefile targets can also be used for dev convenience

**Build image using minikube docker daemon, create k8s secrets, apply k8s manifests**
```
make k8s-deploy
```

**Open k8s pmt service tunnel to our localhost**
```
make k8s-tunnel-pmt
```

**Open k8s db tunnel to our localhost**
```
make k8s-tunnel-db
```

**Unapply k8s manifests, delete k8s secrets, delete docker image**
```
make k8s-down
```

## Detailed k8s deployment instructions

**Build docker image**

For production, we would need to publish our image to a container registry.

For our minikube cluster to access the image, let the shell use the minikube docker daemon before building the image
(this is faster than an alternative: `minikube image load pmt-api:local` after the image has been built):
```
eval $(minikube docker-env)
```

Build the image:

```
docker build -t pmt-api:local .
```

Or to clean up later:

```
docker rmi pmt-api:local
```

**Create secrets in cluster used by our service**

Create secrets:

```
kubectl create secret generic db-secret \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres \
  --from-literal=POSTGRES_DB=pmt_db
```

Inspect secrets:

```
kubectl get secret db-secret
```

**Apply k8s manifests defined in `k8s/`**

Apply manifests:

```
kubectl apply -f k8s/
```

Check pods:

```
kubectl get pods
```

Check logs:

```
kubectl logs service/pmt
kubectl logs service/db
```

**Port forward services to localhost**

```
kubectl port-forward service/pmt 8080:8080 9090:9090
```

Then see and run the commands documented in the "# Services" section of this README.md

**Teardown our cluster specfications**

```
kubectl delete -f k8s/
kubectl delete secret db-secret
```
