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

# Verificar se o RabbitMQ está rodando
$rabbitmqRunning = minikube kubectl get pods -l app=rabbitmq | Select-String "rabbitmq"
if ($rabbitmqRunning) {
    Write-Host "RabbitMQ já está rodando. Pulando a criação do RabbitMQ."
} else {
    Write-Host "RabbitMQ não encontrado. Criando RabbitMQ..."
    minikube kubectl -- apply -f "kubernetes/rabbitmq-cluster.yaml"
}

# Verificar se o Minio está rodando
$minioRunning = minikube kubectl get pods -l app=minio | Select-String "minio"
if ($minioRunning) {
    Write-Host "Minio já está rodando. Pulando a criação do Minio."
} else {
    Write-Host "Minio não encontrado. Criando Minio..."
    minikube kubectl -- apply -f "kubernetes/minio-statefulset.yaml"
}

# Removendo os serviços frontend, backend, e worker antes de recriar
Write-Host "Removendo serviços existentes (frontend, backend, worker)..."
$servicesToRemove = @("kubernetes/backend-deployment.yaml", "kubernetes/frontend-deployment.yaml", "kubernetes/job-creator.yaml")
foreach ($file in $servicesToRemove) {
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

# Aplicando novamente os manifests do Kubernetes para os serviços frontend, backend e worker
Write-Host "Recriando serviços no Kubernetes..."
foreach ($file in $servicesToRemove) {
    Write-Host "Aplicando $file"
    minikube kubectl -- apply -f $file
}

Write-Host "Deploy finalizado com sucesso!"

