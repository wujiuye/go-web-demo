FROM --platform=linux/amd64 docker.io/golang:1.19 as builder
WORKDIR /workspace
COPY go.mod go.mod
#COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web-demo main.go

FROM --platform=linux/amd64 docker.io/golang:1.19
WORKDIR /usr/local/lizhi/web-demo
COPY --from=builder /workspace/web-demo .
RUN apt-get update && apt-get install -y curl
ENTRYPOINT ["./web-demo"]
