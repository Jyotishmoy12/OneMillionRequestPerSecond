param(
    [string]$Url = "http://host.docker.internal:8080/v1/items/1",
    [int[]]$ConcurrencyLevels = @(100, 250, 500),
    [string]$Duration = "30s"
)

Write-Host "Local benchmark started"
Write-Host "Target URL: $Url"
Write-Host "Duration: $Duration"
Write-Host ""

foreach ($Concurrency in $ConcurrencyLevels) {
    Write-Host "========================================"
    Write-Host "Concurrency: $Concurrency"
    Write-Host "Duration: $Duration"
    Write-Host "URL: $Url"
    Write-Host "========================================"

    docker run --rm alpine/bombardier `
        -c $Concurrency `
        -d $Duration `
        -l $Url

    Write-Host ""
}