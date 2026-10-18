# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.27.2-alpine AS build

WORKDIR /src

# download modules first so this layer is cached when only code changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# static binary, no cgo, so it runs on a distroless image with no libc
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w \
      -X github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version.Version=${VERSION} \
      -X github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version.Commit=${COMMIT} \
      -X github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version.BuildTime=${BUILD_TIME}" \
    -o /out/server ./cmd/server

# ---- runtime stage ----
# distroless: no shell, no package manager, runs as uid 65532 (nonroot)
FROM gcr.io/distroless/static-debian12:nonroot

ARG VERSION=dev
ARG COMMIT=unknown
LABEL org.opencontainers.image.source="https://github.com/BHAV0207/SuperLeap-BackendAssignment" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}"

COPY --from=build /out/server /server

USER nonroot:nonroot
EXPOSE 8080
ENV GIN_MODE=release

ENTRYPOINT ["/server"]
