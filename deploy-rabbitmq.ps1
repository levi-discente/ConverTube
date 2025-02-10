# deploy-rabbitmq.ps1
Write-Host "Deploy do RabbitMQ iniciado..."

# Deleta o serviço antes de recriar
Write-Host "Deletando RabbitMQ..."
minikube kubectl -- delete -f kubernetes/rabbitmq-cluster.yaml --ignore-not-found

# Recria o serviço
Write-Host "Aplicando configuração do RabbitMQ..."
minikube kubectl -- apply -f kubernetes/rabbitmq-cluster.yaml

Write-Host "Deploy do RabbitMQ finalizado com sucesso!"

