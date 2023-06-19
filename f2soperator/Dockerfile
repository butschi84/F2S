# Use the official Go image as the base image
FROM ubuntu:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# install go
RUN apt update && apt install golang -y

# Build the Go project
RUN go build -o f2s

# expose api
EXPOSE 8080

# expose metrics
EXPOSE 8081

ENV http_proxy      ""
ENV https_proxy     ""

# Set the command to run when the container starts
CMD ["./f2s"]
