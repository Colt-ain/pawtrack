FROM golang:1.23-alpine AS build
WORKDIR /app
RUN apk add --no-cache build-base git
COPY go.mod ./
RUN go mod download
COPY . .
ENV GOTOOLCHAIN=auto
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go --output ./docs
RUN go build -o pawtrack ./main.go

FROM alpine:3.20
WORKDIR /srv
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /app/pawtrack /usr/local/bin/pawtrack
COPY migrations ./migrations
COPY docs /srv/docs
RUN ls -R /srv
EXPOSE 8080
ENV ADDR=:8080
ENTRYPOINT ["/usr/local/bin/pawtrack"]
