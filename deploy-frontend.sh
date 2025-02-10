#!/bin/bash
set -e

echo "Configurando ambiente Docker para Minikube..."
eval $(minikube docker-env)

echo "Construindo a imagem do Frontend..."
cd frontend
docker build -t frontend .
cd ..

echo "Realizando deploy do Frontend no Kubernetes..."
cd kubenertes
minikube kubectl -- apply -f frontend-deployment.yaml
cd ..

echo "Deploy do Frontend realizado com sucesso!"
