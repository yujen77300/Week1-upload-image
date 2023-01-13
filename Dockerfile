# # builde stage
# FROM golang:1.19-alpine3.16 AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o main .

# # run stage
# FROM alpine:3.16
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY config.toml .

# EXPOSE 808
# CMD ["./main"]


FROM golang:1.19-alpine3.16
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]