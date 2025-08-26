# Sample Dockerfile for SBOM generation and vulnerability scanning
FROM node:18-alpine AS base

# Install dependencies
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# Add a vulnerable package for demo purposes
FROM base AS development
RUN npm install lodash@4.17.20

# Create a simple Express app
COPY . .
EXPOSE 3000

# Run as non-root user for security
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

CMD ["npm", "start"]