FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
