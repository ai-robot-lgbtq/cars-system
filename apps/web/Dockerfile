# syntax=docker/dockerfile:1.6
FROM node:20-alpine
WORKDIR /app
# P0 占位 - 后续 Task 12 添加 pnpm install / build
RUN echo "console.log('web placeholder')" > index.js
EXPOSE 3000
CMD ["node", "index.js"]
