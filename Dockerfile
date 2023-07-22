
# build f2soperator
FROM golang:1.19 AS f2soperator
WORKDIR /app
COPY ./f2soperator .
RUN go build -o f2s

# build f2sfrontend
FROM node:18.12.1-alpine as f2sfrontend
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY ./f2sfrontend/package.json ./
COPY ./f2sfrontend/package-lock.json ./
RUN npm ci 
RUN npm install react-scripts@3.4.3 -g --silent
COPY ./f2sfrontend .
RUN npm run build

# put together all files
FROM golang:1.19
WORKDIR /app
COPY --from=f2soperator /app/f2s ./f2s
COPY ./static ./static
COPY --from=f2sfrontend /app/build/ ./static/frontend/

# expose api
EXPOSE 8080

# expose metrics
EXPOSE 8081

ENV http_proxy      ""
ENV https_proxy     ""

# Set the command to run when the container starts
CMD ["/app/f2s"]