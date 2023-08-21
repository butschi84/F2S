# ###################################
# f2s fizzlet base image
# ###################################
FROM romanhuesler/f2sfizzlet AS f2sfizzlet

# ###################################
# build nodejs project
# ###################################
FROM node:18.12.1

# Create app directory
WORKDIR /usr/src/app

# Install dependencies (npm is already installed)
COPY ./executables/ ./
RUN npm install
COPY --from=f2sfizzlet /app/f2sfizzlet ./f2sfizzlet

EXPOSE 9092
CMD ["./init.sh"]