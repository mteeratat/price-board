# Stage 1: build
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /out/fetcher ./cmd/fetcher

# Stage 2: run (no shell, no compiler, non-root)
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/ /app/
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
