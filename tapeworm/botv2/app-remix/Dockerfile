FROM node:17.9.0-slim

WORKDIR /workspace/app

ADD package.json package-lock.json .
RUN npm install

ADD . .
RUN npm run build

CMD ["npm", "run", "start"]