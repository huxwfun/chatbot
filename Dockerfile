FROM golang:1.20-alpine

WORKDIR /app
COPY go.mod go.sum README.md ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY frontend/out ./frontend/out

RUN echo $(ls -l) --no-cache
RUN echo $(go env) --no-cache
RUN go build -o /app/docker-chatbot cmd/main.go

EXPOSE 8080

# Run
CMD ["./docker-chatbot"]

