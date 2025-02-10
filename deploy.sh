#!/bin/bash
set -e # Encerra o script se algum comando retornar erro

# 1. Configura o ambiente Docker para usar o daemon do Minikube
echo "Configurando ambiente Docker para Minikube..."
eval $(minikube docker-env)

# 2. Build dos serviços

## Build do Backend
echo "Construindo a imagem do backend..."
cd backend
docker build -t backend .
cd ..

## Build do Frontend
echo "Construindo a imagem do frontend..."
cd frontend
docker build -t frontend .
cd ..

## Build do Worker
echo "Construindo a imagem do worker..."
cd Worker
docker build -t worker_v1 .
cd ..

# 3. Verifica se os serviços RabbitMQ e Minio já estão criados antes de recriar

echo "Verificando se RabbitMQ e Minio já estão criados..."

# Verifica se o RabbitMQ já está rodando
if ! minikube kubectl get pods -l app=rabbitmq | grep -q "rabbitmq"; then
  echo "RabbitMQ não encontrado. Criando RabbitMQ..."
  minikube kubectl -- apply -f rabbitmq-cluster.yaml
else
  echo "RabbitMQ já está rodando. Pulando a criação."
fi

# Verifica se o Minio já está rodando
if ! minikube kubectl get pods -l app=minio | grep -q "minio"; then
  echo "Minio não encontrado. Criando Minio..."
  minikube kubectl -- apply -f minio-statefulset.yaml
else
  echo "Minio já está rodando. Pulando a criação."
fi

# 4. Deploy dos recursos no Kubernetes

echo "Realizando deploy dos recursos no Kubernetes..."
cd kubenertes

# Aplica os manifests de frontend, backend e worker
minikube kubectl -- delete -f backend-deployment.yaml --ignore-not-found
minikube kubectl -- delete -f frontend-deployment.yaml --ignore-not-found
minikube kubectl -- delete -f job-creator.yaml --ignore-not-found
minikube kubectl -- apply -f backend-deployment.yaml
minikube kubectl -- apply -f frontend-deployment.yaml
minikube kubectl -- apply -f job-creator.yaml

echo "Deploy realizado com sucesso!"
cho "Deploy realizado com sucesso!"
