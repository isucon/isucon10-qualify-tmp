FROM node:12.17

EXPOSE 3000

WORKDIR /frontend
COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci
COPY . .

RUN npm run build

CMD ["npm", "run", "export"]
