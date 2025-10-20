FROM golang:1.24 AS builder

WORKDIR /complexity-analyzer

RUN apt-get update && apt-get install -y gcc libc6-dev make git

ENV CGO_ENABLED=1

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o complexity-analyzer

FROM golang:1.24

WORKDIR /complexity-analyzer

COPY --from=builder /complexity-analyzer/complexity-analyzer .
COPY --from=builder /complexity-analyzer/templates ./templates
COPY --from=builder /complexity-analyzer/static ./static

EXPOSE 3000

CMD ["./complexity-analyzer"]