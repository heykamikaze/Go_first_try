FROM golang:alpine

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
COPY static ./
COPY templates ./
COPY main.go ./
RUN go mod download && go mod verify

EXPOSE 8080:8080

CMD ["go", "run", "main.go"]
