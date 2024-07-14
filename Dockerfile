FROM golang:latest as builder

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/api-server ./cmd/api-server


FROM scratch
COPY --from=builder /app/bin/api-server /api-server
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD [ "/api-server" ]


