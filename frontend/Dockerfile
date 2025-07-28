# build stage
FROM oven/bun:1 AS build
WORKDIR /app
COPY package.json bun.lockb ./
RUN bun install
COPY src/ ./src
COPY index.html ./

RUN bun run build

# final stage
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx config for Docker Compose (default)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create nginx config for Kubernetes
RUN echo 'server {' > /etc/nginx/conf.d/default-k8s.conf && \
    echo '    listen 80;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    server_name localhost;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    root /usr/share/nginx/html;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    index index.html;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    # Handle API requests by proxying to backend' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    location /api/ {' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        proxy_pass http://react-go-app-backend:8080;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        proxy_set_header Host $host;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        proxy_set_header X-Forwarded-Proto $scheme;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    }' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    # Handle all other requests by serving the React app' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    location / {' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '        try_files $uri $uri/ /index.html;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    }' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    # Optional: Add some basic security headers' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    add_header X-Frame-Options "SAMEORIGIN" always;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    add_header X-Content-Type-Options "nosniff" always;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '    add_header X-XSS-Protection "1; mode=block" always;' >> /etc/nginx/conf.d/default-k8s.conf && \
    echo '}' >> /etc/nginx/conf.d/default-k8s.conf

# Create entrypoint script to choose the right config based on environment
RUN echo '#!/bin/sh' > /docker-entrypoint.sh && \
    echo 'if [ "$KUBERNETES_SERVICE_HOST" ]; then' >> /docker-entrypoint.sh && \
    echo '    echo "Running in Kubernetes, using k8s nginx config"' >> /docker-entrypoint.sh && \
    echo '    cp /etc/nginx/conf.d/default-k8s.conf /etc/nginx/conf.d/default.conf' >> /docker-entrypoint.sh && \
    echo 'else' >> /docker-entrypoint.sh && \
    echo '    echo "Running in Docker Compose, using default nginx config"' >> /docker-entrypoint.sh && \
    echo '    # Remove the k8s config to avoid nginx validation errors' >> /docker-entrypoint.sh && \
    echo '    rm -f /etc/nginx/conf.d/default-k8s.conf' >> /docker-entrypoint.sh && \
    echo 'fi' >> /docker-entrypoint.sh && \
    echo 'nginx -g "daemon off;"' >> /docker-entrypoint.sh && \
    chmod +x /docker-entrypoint.sh

# expose 80 by default
EXPOSE 80

# Use our custom entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"] 