FROM node:19-alpine as base

WORKDIR /usr/src/frontend

COPY package.json .
COPY yarn.lock .

RUN yarn

COPY . .

FROM base as builder

RUN yarn build

FROM nginx:1.23-alpine

COPY --from=builder /usr/src/frontend/dist/ /usr/share/nginx/html

COPY --from=builder /usr/src/frontend/nginx.conf /etc/nginx/conf.d/default.conf
