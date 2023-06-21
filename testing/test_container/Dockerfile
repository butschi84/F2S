FROM node:18.12.1

# Create app directory
WORKDIR /usr/src/app

# Install dependencies (npm is already installed)
COPY package*.json ./
RUN npm install

# Bundle app source
COPY . .

EXPOSE 3000

CMD [ "npm", "start" ]