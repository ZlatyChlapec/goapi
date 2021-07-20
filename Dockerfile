FROM golang:1.16.6-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/api

# production image
FROM alpine

COPY --from=build /app/api /app/api

EXPOSE 8080
CMD ["/app/api"]