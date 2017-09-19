var ws = new WebSocket("ws://localhost:9000/ws");

var App = new Vue({
    delimiters: ['${', '}'],

    el: '#app',

    data: {
        outlets: [],
        serverDT: '',
        test: 'hello',
        UtilPercent: 0,
        ServerT: '',
        operatorName: '',
        terminalName: ''
    },

    created: function() {
        this.operatorName = document.getElementById("operator").value;
        console.log(this.operatorName);

        ws.onopen = function(evt) {
            console.log('Connected to ws.');
            console.log('joining room.');
            ws.send(JSON.stringify({
                "Event": "JOIN",
                "Acct": "user1",
                "Outlet": this.operatorName
            }));
        }.bind(this);

        ws.onmessage = function (evt) {
            var t = JSON.parse(evt.data);
            if (t.Event === "OUTLET-UPDATE") {
                for (var k in t.Outsum) {
                    if (k === this.operatorName) {
                        this.outlets = t.Outsum[k];
                    }
                }
            } else if (t.Event === "TIME-UPDATE") {
                this.ServerT = t.ServerT;
            }
        }.bind(this);

        ws.onerror = function (evt) {
            console.log('Error ' + evt.data);
        }.bind(this);

        ws.onclose = function (evt) {
            console.log('Disconnected to ws.');
        }.bind(this);
    },

    methods: {
        getYear: function () {
            var d = new Date();
            return d.getFullYear();
        },

        getInfoBoxClass: function(t, o) {
            var c;
            var perUtil = ((o / t) * 100).toFixed(2);
            if (perUtil > 75) {
                // c = ['info-box', 'bg-green'];
                c = {'info-box': true, 'bg-green': true};
            }
            else if (perUtil > 50) {
                c = {'info-box': true, 'bg-aqua': true};
            }
            else if (perUtil > 25) {
                c = {'info-box': true, 'bg-yellow': true};
            }
            else {
                c = {'info-box': true, 'bg-red': true};
            }
            return c;
        },

        getUtilization: function (t, o) {
            return ((o / t) * 100).toFixed(2);
        }
    }

});