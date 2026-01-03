# Build app
FROM node:lts-alpine3.17 as node
WORKDIR .
RUN mkdir desktop
WORKDIR desktop
COPY . .
RUN npm i -g @angular/cli@17.3.3
RUN npm i
RUN ng build

# Deploy via nginx
FROM nginx:alpine
COPY --from=node /desktop/dist/receipt-wrangler /usr/share/nginx/html
COPY --from=node /desktop/docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
