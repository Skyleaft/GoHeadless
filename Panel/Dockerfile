# Build stage
FROM node:24-alpine AS builder

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml* .npmrc* ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build the application (produces static files in build/ directory)
RUN pnpm build

# Production stage - Nginx will serve static files
# Note: The final serving is handled by the main nginx container
FROM node:24-alpine AS runner

WORKDIR /app

# Copy built files
COPY --from=builder /app/build ./build

# Expose port for preview (not used directly in production)
EXPOSE 3000

# For static serving, we just need the output
CMD ["echo", "Build complete - files in /app/build"]