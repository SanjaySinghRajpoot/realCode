FROM node:18-alpine

WORKDIR /react-vite-app

RUN apk update && apk add curl

EXPOSE 3000

COPY package.json package-lock.json ./

RUN npm install --silent

COPY . ./

CMD ["npm", "run", "dev"]