FROM node:22.12.0-alpine AS builder

# Define diretório de trabalho
WORKDIR /app

# Copia apenas arquivos essenciais
COPY package.json package-lock.json ./

# Instala dependências sem devDependencies
RUN npm ci --omit=dev

# Copia o restante do código
COPY . .

# Gera o build da aplicação
RUN npm run build

# Etapa final (imagem leve com apenas o necessário)
FROM node:22.12.0-alpine

# Define diretório de trabalho
WORKDIR /app

# Copia apenas os arquivos necessários da etapa anterior
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Expor porta
EXPOSE 3000

# Iniciar Next.js no modo de produção
CMD ["npm", "run", "start"]

