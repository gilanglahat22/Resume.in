FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the application code
COPY . .

# Build the application for production
RUN npm run build -- --configuration production

# Check the build output directory structure
RUN ls -la /app/dist

# Use Nginx to serve the application
FROM nginx:1.21-alpine

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy the Angular app build from the builder stage
# For Angular Universal apps, we need the browser directory
COPY --from=builder /app/dist/client/browser /usr/share/nginx/html

# Copy our custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]