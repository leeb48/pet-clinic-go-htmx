FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o pet-clinic

EXPOSE 7771

CMD [ "/app/cmd/pet-clinic" ]