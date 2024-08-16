$env:PRODUCT_NAME="wcb"
Write-Output("Building $env:PRODUCT_NAME`n")

Write-Output("Current GOOS $env:GOOS : Current GOARCH :$env:GOARCH.`n")
$env:OLDGOOS=$env:GOOS
$env:OLDGOARCH=$env:GOARCH

$env:GOARCH="amd64"
$env:GOOS="linux"
Write-Output("Linux x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -ldflags="-w -s" -o distr/$env:PRODUCT_NAME
Write-Output("Linux build done.`n")

$env:GOOS=$env:OLDGOOS
$env:GOARCH=$env:OLDGOARCH