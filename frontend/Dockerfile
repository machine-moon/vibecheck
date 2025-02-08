FROM node:22.13-alpine
WORKDIR /src
ENV NODE_ENV=production
COPY package*.json ./
RUN npm install -g yarn && yarn install --frozen-lockfile || yarn install
COPY --chown=node:node . ./
USER node

EXPOSE 3000
COPY secrets.yml .
RUN source secrets.yml
CMD ["yarn", "deploy"]
