# deploy-backend.ps1
Write-Host "Deploy do Backend iniciado..."

# Deleta o serviço antes de recriar
Write-Host "Deletando backend..."
minikube kubectl -- delete -f kubernetes/backend-deployment.yaml --ignore-not-found

# Build da imagem no Minikube
Write-Host "Construindo a imagem do backend..."
Set-Location backend
minikube image build -t "backend:latest" .
Set-Location ..

# Recria o serviço
Write-Host "Aplicando configuração do backend..."
minikube kubectl -- apply -f kubernetes/backend-deployment.yaml

Write-Host "Deploy do Backend finalizado com sucesso!"

