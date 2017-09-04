var ws = new WebSocket("ws://localhost:9000/ws");

var App = new Vue({
    delimiters: ['${', '}'],

    el: '#app',

    data: {
        franchise: [],
        serverDT: '',
        test: 'hello',
        UtilPercent: 0,
        ServerT: ''
    },

    created: function() {
        ws.onopen = function(evt) {
            console.log('Connected to ws.');
        }.bind(this);

        ws.onmessage = function (evt) {
            var t = JSON.parse(evt.data);
            if (t.Event === "DB-UPDATE") {
                this.franchise = t.Netsum;
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