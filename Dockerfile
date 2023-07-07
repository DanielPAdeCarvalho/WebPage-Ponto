############################
# STEP 1 build executable binary
############################
FROM golang:1.20.5-alpine3.18 AS builder

WORKDIR /app

# Fetch dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source and build the application.
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o  main .

############################
# STEP 2 build a small image
############################
FROM scratch

# Create appuser
USER 1001

COPY --from=builder /app/main /app/main
COPY --from=builder /app/templates /templates
COPY --from=builder /app/assets /assets

EXPOSE 8080
ENTRYPOINT [ "/app/main" ]
