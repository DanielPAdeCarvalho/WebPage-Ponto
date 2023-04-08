############################
# STEP 1 build executable binary
############################
FROM golang:1.20.3-alpine3.17 as builder
WORKDIR /app
COPY . .
WORKDIR /app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o  main .

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /app/main /app/main
COPY --from=builder /app/templates /templates
COPY --from=builder /app/assets /assets

EXPOSE 8080
ENTRYPOINT [ "/app/main" ]
CMD ["/app/main"]