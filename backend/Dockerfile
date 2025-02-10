FROM node:22.12.0-alpine AS builder

# Define diretório de trabalho
WORKDIR /app

# Copia apenas arquivos essenciais
COPY package.json package-lock.json ./

# Instala dependências sem instalar devDependencies
RUN npm ci --omit=dev

# Copia o restante do código
COPY . .

# Compila a aplicação (gera dist/)
RUN npm run build

# Etapa final (imagem leve)
FROM node:22.12.0-alpine

# Define diretório de trabalho
WORKDIR /app

# Copia apenas os arquivos necessários da etapa anterior
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Expor porta (caso necessário)
EXPOSE 3000

# Comando de inicialização
CMD ["node", "dist/main"]

