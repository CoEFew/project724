# Dockerfile (ฝั่ง backend)
FROM golang:1.23.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /server ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=build /server /server
ENV PORT=8080
EXPOSE 8080
CMD ["/server"]
