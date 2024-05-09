"package utils

var BuildTime=`"$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`"" | Out-File -FilePath .\utils\version.go -Encoding ascii