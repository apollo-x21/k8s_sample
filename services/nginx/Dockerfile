FROM node:18-alpine AS build
WORKDIR /frontend
COPY services/frontend/package.json services/frontend/tsconfig.json services/frontend/.umirc.ts ./
COPY services/frontend/src ./src
RUN npm install && npm run build

FROM nginx:1.25-alpine
COPY services/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /frontend/dist /usr/share/nginx/html
EXPOSE 80
