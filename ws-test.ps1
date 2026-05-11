$tokens = @()

1..40 | ForEach-Object {

    $login = Invoke-RestMethod -Method POST `
    -Uri "http://localhost:8080/auth/login" `
    -ContentType "application/json" `
    -Body "{`"username`":`"testuser$_`",`"password`":`"123456`"}"

    $tokens += $login.token
}

$i = 1

foreach ($token in $tokens) {

    Start-Job {
        param($token, $i)

        Add-Type -AssemblyName System.Net.WebSockets

        $ws = [System.Net.WebSockets.ClientWebSocket]::new()

        $uri = [Uri]"ws://127.0.0.1:8080/ws?token=$token&manga_id=one-piece"

        $ws.ConnectAsync(
            $uri,
            [Threading.CancellationToken]::None
        ).Wait()

        Write-Host "User $i connected"

        Start-Sleep -Seconds 60

        $ws.Dispose()

    } -ArgumentList $token, $i

    $i++
}