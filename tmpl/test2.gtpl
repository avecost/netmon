<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>test</title>
</head>
<body>

<script>
    var ws = new WebSocket("ws://localhost:9000/ws");
    ws.onopen = function()
    {
        console.log("we are connected");
        console.log("sending join");
        ws.send(JSON.stringify({
            "Event": "JOIN",
            "Acct": "user2",
            "Outlet": "TGXI"
        }));
    };

    ws.onmessage = function(evt)
    {
        console.log(evt)
    };

    ws.onclose = function()
    {
        console.log("ws closed")
    };

</script>
</body>
</html>