FROM node:latest AS builder
ADD . /app
WORKDIR /app
RUN npm install
RUN npm run build-prod

# final stage
FROM node:latest
EXPOSE 4001
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
ENTRYPOINT ["npm", "start"]