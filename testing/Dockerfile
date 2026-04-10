FROM mcr.microsoft.com/playwright:v1.58.2-jammy

WORKDIR /app

COPY testing/package.json ./
RUN npm install

COPY testing/ ./
COPY services/db/init.sql ./services/db/init.sql

CMD ["npm", "test"]
