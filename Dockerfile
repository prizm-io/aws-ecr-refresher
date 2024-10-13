FROM golang:1.22.2 as builder

WORKDIR /app

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 go build -o /app/application /app/*.go

FROM scratch

COPY --from=builder /app/application /application

ENTRYPOINT ["/application"]
