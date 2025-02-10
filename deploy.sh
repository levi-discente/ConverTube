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

# 3. Deploy dos recursos no Kubernetes

echo "Realizando deploy dos recursos no Kubernetes..."
cd kubenertes

# Opcional: remover recursos antigos (usa --ignore-not-found para não dar erro se não existirem)
minikube kubectl -- delete -f rabbitmq-cluster.yaml --ignore-not-found
minikube kubectl -- delete -f minio-statefulset.yaml --ignore-not-found
minikube kubectl -- delete -f backend-deployment.yaml --ignore-not-found
minikube kubectl -- delete -f frontend-deployment.yaml --ignore-not-found
minikube kubectl -- delete -f job-creator.yaml --ignore-not-found

# Aplica os manifests em uma ordem que faça sentido (dependendo do seu setup)
minikube kubectl -- apply -f rabbitmq-cluster.yaml
minikube kubectl -- apply -f minio-statefulset.yaml
minikube kubectl -- apply -f backend-deployment.yaml
minikube kubectl -- apply -f frontend-deployment.yaml
minikube kubectl -- apply -f job-creator.yaml

echo "Deploy realizado com sucesso!"
