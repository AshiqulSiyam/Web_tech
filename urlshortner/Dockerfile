FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 3000
RUN go build -o server

CMD [ "./server" ]
