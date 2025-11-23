# Build stage
FROM golang:1.24.4-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN rm -rf volumes
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

EXPOSE 8088
CMD ["./main"]