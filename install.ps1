# Esse script precisa ser executado com privilégios administrativos

# Verifica se o script está rodando como Administrador.
$windowsIdentity = [Security.Principal.WindowsIdentity]::GetCurrent()
$windowsPrincipal = New-Object Security.Principal.WindowsPrincipal($windowsIdentity)
if (-not $windowsPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "Erro: Execute esse script como Administrador (run as administrator)."
    exit 1
}

Write-Host "Script executando com privilégios de Administrador..."

# 1. Instalar plugins via Minikube

Write-Host "Instalando o RabbitMQ Cluster Operator..."
minikube kubectl -- apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml

Write-Host "Instalando o Ingress NGINX..."
minikube kubectl -- apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

# 2. Obter o IP do Minikube
$minikubeIP = minikube ip
if ([string]::IsNullOrEmpty($minikubeIP)) {
    Write-Error "Erro: Não foi possível obter o IP do Minikube."
    exit 1
}
Write-Host "Minikube IP: $minikubeIP"

# 3. Atualizar o arquivo de hosts
$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
$domains = @("frontend.local", "minio.local", "rabbitmq.local", "backend.local")

# Lê o conteúdo atual do arquivo de hosts
$hostsContent = Get-Content $hostsPath

# Remove entradas antigas para os domínios (usando regex para identificar linhas que terminem com o domínio)
foreach ($domain in $domains) {
    $hostsContent = $hostsContent | Where-Object { $_ -notmatch "\s+$domain(\s|$)" }
}

# Gera as novas entradas
$newEntries = @()
foreach ($domain in $domains) {
    $newEntries += "$minikubeIP $domain"
}

# Escreve o conteúdo atualizado no arquivo de hosts
$updatedContent = $hostsContent + $newEntries
$updatedContent | Set-Content $hostsPath -Force

Write-Host "Arquivo hosts atualizado:"
Get-Content $hostsPath | Select-String -Pattern "frontend.local|minio.local|rabbitmq.local|backend.local"

# 4. Executar o deploy geral
if (Test-Path ".\deploy.ps1") {
    Write-Host "Executando o deploy geral (deploy.ps1)..."
    .\deploy.ps1
} else {
    Write-Error "Arquivo deploy.ps1 não encontrado."
    exit 1
}

Write-Host "Instalação e deploy concluídos com sucesso!"

