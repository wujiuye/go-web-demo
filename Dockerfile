FROM public.ecr.aws/docker/library/golang:1.19 AS builder
WORKDIR /workspace
COPY go.mod go.mod
RUN go mod download
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web-demo main.go

FROM public.ecr.aws/docker/library/golang:1.19
WORKDIR /usr/local/app/web-demo
COPY --from=builder /workspace/web-demo .
RUN apt-get update && apt-get install -y curl
ENTRYPOINT ["./web-demo"]
