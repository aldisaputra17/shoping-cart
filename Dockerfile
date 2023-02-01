FROM golang:alpine as builder

LABEL maintainer="Aldi Saputra Dalimunthe <aldisaput17@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .
RUN go build -o main

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]
