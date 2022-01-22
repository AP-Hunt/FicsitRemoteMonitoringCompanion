$template = Get-Content -Path "./readme.tpl.md"
$table = ../Companion/bin/Companion.exe -ShowMetrics

$template -replace "__METRICS_TABLE__",$table