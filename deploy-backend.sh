#!/bin/bash
set -e

echo "Configurando ambiente Docker para Minikube..."
eval $(minikube docker-env)

echo "Construindo a imagem do Backend..."
cd backend
docker build -t backend .
cd ..

echo "Realizando deploy do Backend no Kubernetes..."
cd kubenertes
minikube kubectl -- apply -f backend-deployment.yaml
cd ..

echo "Deploy do Backend realizado com sucesso!"
