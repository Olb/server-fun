# --- Build React Frontend ---
    FROM node:20-alpine AS builder
    WORKDIR /web
    
    COPY package.json package-lock.json ./
    RUN npm install
    
    COPY . ./
    RUN npm run build
    
    # --- Serve with Nginx ---
    FROM nginx:alpine
    WORKDIR /usr/share/nginx/html
    
    # Copy built React app
    COPY --from=builder /web/dist .
    
    # Expose port
    EXPOSE 80
    CMD ["nginx", "-g", "daemon off;"]
    