# 1. Aşama: Build
FROM node:22.16.0-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# 2. Aşama: Production (Nginx)
FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf

# Vite varsayılan olarak 'dist' klasörüne çıktı verir. 
# (Eğer Create React App kullanıyorsan 'dist' yerine 'build' yazmalısın)
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]