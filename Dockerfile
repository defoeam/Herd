# syntax=docker/dockerfile:1

# --- BUILD STAGE ---
FROM golang:1.21 AS build-stage

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
RUN CGO_ENABLED=0 GOOS=linux go build -o /herd

# --- DEPLOY STAGE ---
# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

# Copy the log directory from build stage
COPY --from=build-stage /app/log /app/log

COPY --from=build-stage /herd /herd

EXPOSE 8080

USER root:root

ENTRYPOINT ["/herd"]
