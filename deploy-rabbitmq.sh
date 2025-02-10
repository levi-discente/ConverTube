#!/bin/bash
set -e

echo "Realizando deploy do RabbitMQ no Kubernetes..."
cd kubenertes
minikube kubectl -- apply -f rabbitmq-cluster.yaml
cd ..

echo "Deploy do RabbitMQ realizado com sucesso!"
