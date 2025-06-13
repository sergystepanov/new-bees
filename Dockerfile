FROM golang:1.24-alpine AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

FROM alpine:3.22

COPY --from=build /usr/local/bin/app ./app

EXPOSE 8081
ENV PORT 8081

ENTRYPOINT ["./app"]