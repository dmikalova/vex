# Build the WebAssembly client and the native server that serves it, then ship
# both on a minimal distroless image. The server honors $PORT (set by Cloud Run).

FROM golang:1.26 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# WebAssembly client bundle served at /web/app.wasm by go-app.
RUN GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/web
# Native server binary.
RUN CGO_ENABLED=0 go build -o /vactrol-web ./cmd/web

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /vactrol-web /vactrol-web
COPY --from=build /app/web ./web

EXPOSE 8080
ENTRYPOINT ["/vactrol-web"]
