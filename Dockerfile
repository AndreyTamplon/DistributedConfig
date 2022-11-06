FROM golang:latest
RUN mkdir /config_service
WORKDIR /config_service

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/config_service/main.go app/
COPY . .


RUN go build -o config_service app/main.go

EXPOSE 8084

CMD [ "bash", "./migrate.sh"]
CMD [ "./config_service" ]


