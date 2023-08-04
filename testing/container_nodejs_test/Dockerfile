# ###################################
# f2s fizzlet base image
# ###################################
FROM romanhuesler/f2sfizzlet AS f2sfizzlet

# ###################################
# build nodejs project
# ###################################
FROM node:18.12.1 as builder

# Create app directory
WORKDIR /usr/src/app

# Install dependencies (npm is already installed)
COPY ./executables/ ./
RUN npm install
RUN npm run build

# ###################################
# definitive image
# ###################################
FROM node:18.12.1 as runner
WORKDIR /usr/src/app
COPY --from=f2sfizzlet /app/f2sfizzlet ./f2sfizzlet
COPY --from=builder /usr/src/app/build ./
COPY ./executables/start.sh ./
CMD ["./f2sfizzlet"]