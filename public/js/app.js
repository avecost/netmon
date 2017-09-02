var ws = new WebSocket("ws://whisky.24bet7.com:9000/ws");

var App = new Vue({
    el: '#app',

    data: {
        franchise: '',
        serverDT: '',
        showLog: '',
        test: 'hello',
        UtilPercent: 0,
        ServerT: ''
    },

    created: function() {
        ws.onopen = function(evt) {
            console.log('Connected to ws.');
        }.bind(this);

        ws.onmessage = function (evt) {
            this.showLog = evt.data;

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
            this.UtilPercent = ((o / t) * 100).toFixed(2);
            if (this.UtilPercent > 75) {
                c = ['info-box', 'bg-green'];
            }
            else if (this.UtilPercent > 50) {
                c = ['info-box', 'bg-aqua'];
            }
            else if (this.UtilPercent > 25) {
                c = ['info-box', 'bg-yellow'];
            }
            else {
                c = ['info-box', 'bg-red'];
            }
            return c;
        }
    }

});