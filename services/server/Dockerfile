FROM golang:1.16.0-buster as builder

WORKDIR /app
COPY . ./
RUN go build -a -o ./bin/svc

# runtime image
FROM debian:buster-slim
RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy our static executable
COPY --from=builder /app/bin/svc /svc

# Port on which the service will be exposed.
EXPOSE 8080 8888 9100

WORKDIR /

# Run the svc binary.
CMD ["./svc"]
