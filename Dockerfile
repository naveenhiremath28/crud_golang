FROM golang:1.25 AS build

WORKDIR /app
COPY . .
RUN go clean --modcache
RUN go mod tidy
RUN GOOS=linux go build -o main cmd/main.go

FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /app
COPY --from=build /app/main .

EXPOSE 3000
CMD ["go", "run", "cmd/main.go"]
