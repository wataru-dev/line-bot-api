FROM golang:latest

WORKDIR /app

COPY . .

RUN go install
RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

