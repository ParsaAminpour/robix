FROM golang:1.17

WORKDIR /app

# Install git and air in a single RUN command
RUN go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD [ "air", "-c", ".air.toml" ]