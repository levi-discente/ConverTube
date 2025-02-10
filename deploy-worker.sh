#!/bin/bash
set -e

echo "Configurando ambiente Docker para Minikube..."
eval $(minikube docker-env)

echo "Construindo a imagem do Worker..."
cd Worker
docker build -t worker_v1 .
cd ..

echo "Realizando deploy do Worker no Kubernetes..."
cd kubenertes
minikube kubectl -- apply -f job-creator.yaml
cd ..

echo "Deploy do Worker realizado com sucesso!"
