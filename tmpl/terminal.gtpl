<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Dashboard</title>

    <!-- Latest bootstrap 3.3.7 compiled and minified CSS -->
    <link rel="stylesheet" href="/public/css/bootstrap/bootstrap.min.css">
    <link rel="stylesheet" href="/public/css/bootstrap/bootstrap-theme.min.css">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">

    <!-- App style -->
    <link rel="stylesheet" href="/public/css/app-outlet.css">
</head>
<body style="height: auto; min-height: 100%;">
<div id="app" class="wrapper" style="height: auto; min-height: 100%;">
    <header class="main-header">
        <nav class="navbar">
            <a class="logo" href="#!">NetMon</a><span>Dashboard<small> Ver 1.0</small></span>
            <div class="navbar-custom-menu">
                <ul class="nav navbar-nav">
                    <li>${ ServerT }</li>
                </ul>
            </div>
        </nav>
    </header>
    <div class="content-wrapper" style="min-height: 916px;">
        <input type="hidden" value="{{ .Name }}" name="operator" id="operator" />
        <section class="content">
            <div class="row">
                <div class="col-lg-3 col-xs-6" v-for="(f, k, idx) in outlets">
                    <div v-bind:class="getInfoBoxClass(f.Terminal, f.Online)">
                        <div class="inner">
                            <div class="row text-center">
                                <h5 style="font-size: 12px; font-weight: 700;">${ f.Name }</h5>
                            </div>
                            <div class="row outlet-container">
                                <div class="col-xs-3 text-center">${ getUtilization(f.Terminal, f.Online) }<sup style="font-size: 10px">%</sup></div>
                                <div class="col-xs-3 text-right">Terminals</div>
                                <div class="col-xs-2 text-center">${ f.Terminal }</div>
                                <div class="col-xs-2 text-right">Online</div>
                                <div class="col-xs-2 text-center">${ f.Online }</div>
                            </div>
                        </div>
                        <a href="#!" class="info-box-footer">More info <i class="fa fa-arrow-circle-right"></i></a>
                    </div>
                </div>
            </div>
        </section>
    </div>
    <footer class="main-footer">
        <div class="pull-right hidden-xs">
            <b>Version</b>
            1.0
        </div>
        <strong>
            Copyright &copy; ${ getYear() } IEST,
        </strong>
        All rights reserved.
    </footer>
</div>

<!-- Latest compiled and minified JavaScript -->
<script src="/public/js/bootstrap/bootstrap.min.js"></script>
<script src="/public/js/vue/vue.js"></script>
<script src="/public/js/outlet.js"></script>

</body>
</html>