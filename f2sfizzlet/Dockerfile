FROM golang:1.19 AS builder

# build f2sfizzlet
WORKDIR /app
COPY . .
RUN go build -o f2sfizzlet

# f2sfizzlet image
FROM scratch
WORKDIR /app
COPY --from=builder /app/f2sfizzlet ./f2sfizzlet
RUN chmod a+x ./f2sfizzlet
EXPOSE 8080