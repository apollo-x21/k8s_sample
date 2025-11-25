FROM node:18-alpine AS build
WORKDIR /frontend
COPY frontend/package.json frontend/tsconfig.json frontend/.umirc.ts ./
COPY frontend/src ./src
RUN npm install && npm run build

FROM nginx:1.25-alpine
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /frontend/dist /usr/share/nginx/html
EXPOSE 80
