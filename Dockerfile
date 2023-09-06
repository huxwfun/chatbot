FROM golang:1.20-alpine

WORKDIR /app
COPY go.mod go.sum README.md ./
COPY cmd ./cmd
COPY internal ./internal
COPY frontend/out ./frontend/out

RUN echo $(ls -l) 
RUN echo $(go env)
# RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN go mod download
RUN go build -o /app/docker-chatbot cmd/main.go

EXPOSE 8080

# Run
CMD ["./docker-chatbot"]

