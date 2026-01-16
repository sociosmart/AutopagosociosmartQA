FROM node:23.11-alpine3.21 AS builder

WORKDIR /app

COPY ./package.json .
COPY ./yarn.lock .

# Installing required packages
RUN apk update 

# Installing with npm
COPY . .

RUN yarn

RUN npx vite build --mode production

COPY nginx.conf ./nginx.conf

FROM nginx:1.27-alpine3.21

COPY --from=builder /app/dist/ /usr/share/nginx/html

COPY --from=builder /app/nginx.conf /etc/nginx/conf.d/default.conf

