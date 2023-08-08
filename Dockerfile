FROM golang:1.20 AS pre

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go  ./
COPY memorydb ./memorydb
COPY model ./model

RUN CGO_ENABLED=0 GOOS=linux go build -o /playerScoreManagement

EXPOSE 8080
CMD ["/playerScoreManagement"]

FROM alpine:latest

WORKDIR /

COPY --from=pre /playerScoreManagement /playerScoreManagement
ENV GIN_MODE=release 
EXPOSE 8080
ENTRYPOINT [ "/playerScoreManagement" ]