# ---- build ----
FROM golang:1.25-alpine AS build

WORKDIR /src

# Dependencies are cached separately from the source, so a code-only change
# does not re-download the module graph.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off produces a static binary, which is what the distroless base needs.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server ./cmd

# ---- run ----
# distroless/static carries the CA bundle needed for TLS to a hosted MySQL,
# and nothing else: no shell, no package manager.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=build /out/server /server

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/server"]
