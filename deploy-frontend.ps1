# deploy-frontend.ps1
Write-Host "Deploy do Frontend iniciado..."

# Deleta o serviço antes de recriar
Write-Host "Deletando frontend..."
minikube kubectl -- delete -f kubernetes/frontend-deployment.yaml --ignore-not-found

# Build da imagem no Minikube
Write-Host "Construindo a imagem do frontend..."
Set-Location frontend
minikube image build -t "frontend:latest" .
Set-Location ..

# Recria o serviço
Write-Host "Aplicando configuração do frontend..."
minikube kubectl -- apply -f kubernetes/frontend-deployment.yaml

Write-Host "Deploy do Frontend finalizado com sucesso!"

