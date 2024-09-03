FROM golang:1.22.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /backend

EXPOSE 7540

CMD ["/backend"] 

