# Use the official Go image as the base image
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go project
RUN go build -o f2s

FROM scratch
WORKDIR /app
COPY --from=builder /app/f2s .
COPY ./static ./static

# expose api
EXPOSE 8080

# expose metrics
EXPOSE 8081

ENV http_proxy      ""
ENV https_proxy     ""

# Set the command to run when the container starts
CMD ["./f2s"]