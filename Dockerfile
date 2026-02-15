# Build stage
FROM golang:1.24 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ll-limiter ./cmd/limiter

# Run stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/ll-limiter /ll-limiter

ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/ll-limiter"]
