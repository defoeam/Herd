# syntax=docker/dockerfile:1

# --- BUILD STAGE ---
FROM golang:1.23 AS build-stage

# Accept build arguments with default values
ARG useLogging=false
ARG useSecurity=false

# Set destination for COPY
WORKDIR /app

# Create log directory during build stage
RUN mkdir -p /app/log

# Download Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /herd ./cmd/main.go

# --- DEPLOY STAGE ---
FROM alpine:3.14 AS build-release-stage

WORKDIR /

# Accept build arguments in deploy stage
ARG useLogging
ARG useSecurity

# Set environment variables to pass to the application
ENV USE_LOGGING=${useLogging}
ENV USE_SECURITY=${useSecurity}

# Copy the log directory from build stage
COPY --from=build-stage /app/log /app/log

# Copy the certs directory from build stage
COPY --from=build-stage /app/certs /certs

# Copy the binary from build stage
COPY --from=build-stage /herd /herd

EXPOSE 7878

USER root:root

ENTRYPOINT ["/bin/sh", "-c", "/herd --useLogging=${USE_LOGGING} --useSecurity=${USE_SECURITY}"]
