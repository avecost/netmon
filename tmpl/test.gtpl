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
        console.log(evt)
        // we expect a filtered broadcast based on Outlet parameter
        // if evt.Event == Outlet
    };

    ws.onclose = function()
    {
        console.log("ws closed")
    };

</script>
</body>
</html>