# Verifica se está rodando como Administrador
$windowsIdentity = [Security.Principal.WindowsIdentity]::GetCurrent()
$windowsPrincipal = New-Object Security.Principal.WindowsPrincipal($windowsIdentity)
if (-not $windowsPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "Erro: Execute esse script como Administrador (run as administrator)."
    exit 1
}

Write-Host "Iniciando Deploy no Minikube..."

# Lista de serviços
$services = @("backend", "frontend", "Worker")

# Diretórios dos arquivos de configuração do Kubernetes
$kubernetesFiles = @(
    "kubernetes/backend-deployment.yaml",
    "kubernetes/frontend-deployment.yaml",
    "kubernetes/job-creator.yaml",
    "kubernetes/minio-statefulset.yaml",
    "kubernetes/rabbitmq-cluster.yaml"
)

# Removendo os serviços antes de recriar
Write-Host "Removendo serviços existentes..."
foreach ($file in $kubernetesFiles) {
    Write-Host "Deletando $file"
    minikube kubectl -- delete -f $file --ignore-not-found
}

# Construção das imagens no Minikube
foreach ($service in $services) {
    Write-Host "Construindo imagem do serviço: $service"
    Set-Location $service
    minikube image build -t "$service:latest" .
    Set-Location ..
}

# Aplicando novamente os manifests do Kubernetes
Write-Host "Recriando serviços no Kubernetes..."
foreach ($file in $kubernetesFiles) {
    Write-Host "Aplicando $file"
    minikube kubectl -- apply -f $file
}

Write-Host "Deploy finalizado com sucesso!"

