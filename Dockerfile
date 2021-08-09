# Name this as builder stage
FROM golang:1.16-alpine AS builder
# Move to working directory /library
WORKDIR /go/src/library
# Copy the code into the container
COPY . .
# ...
RUN go build -o server .

# Runtime stage
FROM golang:1.16-alpine
WORKDIR /app
ENV PORT=1323
ENV DSN=library:password@tcp(mysql.c8ajbiky1mzj.ap-southeast-1.rds.amazonaws.com:3306)/library
ENV GOODREAD_APIKEY=6qVbqOjnzhHws97M5gYYA
ENV ISBNDB_APIKEY=46384_fe2f7c5150354b92aa708a8e0fa46463
# Copy binary from builder stage
COPY --from=builder /go/src/library/server .
EXPOSE 1323
RUN mkdir logs
CMD ["./server"]