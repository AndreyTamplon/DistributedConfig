FROM golang:1.19.3-alpine3.16 AS build
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /config_service cmd/config_service/main.go

FROM alpine:3.16 AS production

COPY --from=build config_service .
COPY --from=build app/app.env .
COPY --from=build app/migrate.sh .
COPY --from=build app/migrations ./migrations
RUN apk add --no-cache bash
RUN ["chmod", "+x", "migrate.sh"]
EXPOSE 8084
CMD [ "./config_service" ]

