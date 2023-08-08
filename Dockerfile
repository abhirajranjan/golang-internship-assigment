FROM golang:1.20

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go  ./
COPY memorydb ./memorydb
COPY model ./model

RUN CGO_ENABLED=0 GOOS=linux go build -o /playerScoreManagement
CMD ["/playerScoreManagement"]
