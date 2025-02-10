#!/bin/bash
set -e

# Verifica se o script está sendo executado como root.
if [ "$EUID" -ne 0 ]; then
  echo "Por favor, execute este script como root."
  exit 1
fi

echo "============================================"
echo " Iniciando instalação e configuração no Minikube"
echo "============================================"

# 1. Instala os plugins via Minikube.
echo "Instalando o RabbitMQ Cluster Operator..."
minikube kubectl -- apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml

echo "Instalando o Ingress NGINX..."
minikube kubectl -- apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

# 2. Captura o IP do Minikube.
MINIKUBE_IP=$(minikube ip)
if [ -z "$MINIKUBE_IP" ]; then
  echo "Erro: Não foi possível obter o IP do Minikube."
  exit 1
fi
echo "Minikube IP: $MINIKUBE_IP"

# 3. Atualiza o arquivo /etc/hosts.
HOSTS_ENTRY="
$MINIKUBE_IP frontend.local
$MINIKUBE_IP minio.local
$MINIKUBE_IP rabbitmq.local
$MINIKUBE_IP backend.local
"

echo "Atualizando /etc/hosts..."

# Remove entradas antigas para os domínios (se existirem).
for host in frontend.local minio.local rabbitmq.local backend.local; do
  sed -i.bak "/[[:space:]]$host$/d" /etc/hosts
done

# Adiciona as novas entradas.
echo "$HOSTS_ENTRY" >>/etc/hosts

echo "Arquivo /etc/hosts atualizado:"
grep -E "frontend\.local|minio\.local|rabbitmq\.local|backend\.local" /etc/hosts

# 4. Executa o deploy geral.
echo "Executando o deploy (deploy.sh)..."
if [ -x "./deploy.sh" ]; then
  ./deploy.sh
else
  echo "Arquivo deploy.sh não encontrado ou sem permissão de execução."
  exit 1
fi

echo "============================================"
echo " Instalação e deploy concluídos com sucesso!"
echo "============================================"
