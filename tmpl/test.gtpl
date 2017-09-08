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
            "Acct": "user1",
            "Outlet": "ABLE"
        }));
    };

    ws.onmessage = function(evt)
    {
        var t = JSON.parse(evt.data);
        if (t.Event === "OUTLET-UPDATE") {
            console.log(t.Outsum);
        } else if (t.Event === "TIME-UPDATE") {
            console.log(t.ServerT);
        }

    };

    ws.onclose = function()
    {
        console.log("ws closed")
    };

</script>
</body>
</html>