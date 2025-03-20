
FROM gcr.io/distroless/static:latest

WORKDIR /app

COPY ./bin/main .
COPY ./templates ./templates

EXPOSE 8080

# Run the binary
CMD ["./main"]
