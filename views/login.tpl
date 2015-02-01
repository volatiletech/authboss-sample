<html>
<head>
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css" rel="stylesheet" />
    <link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" rel="stylesheet" />
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js" type="text/javascript"></script>
</head>
<body>
<div class="container-fluid">
    {{if .FlashSuccess}}
    <div class="row">
        <div class="col-xs-offset-3 col-md-6">
            <div class="alert alert-success alert-dismissable" style="margin-top: 75px;">
                <button type="button" class="close" data-dismiss="alert"><span>&times;</span></button>
                {{print .FlashSuccess}}
            </div>
        </div>
    </div>
    {{end}}
    <div class="row" style="margin-top: {{if .FlashSuccess}}25{{else}}75{{end}}px;">
        <div class="col-md-offset-4 col-md-4">
            <div class="panel panel-default">
                <div class="panel-heading">Login</div>
                <div class="panel-body">
                    <form action="/login" method="POST">
                        <div class="form-group{{if .Error}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-user"></i></span>
                                <input type="text" class="form-control" name="username" placeholder="Username" value="{{.Username}}">
                            </div>
                        </div>
                        <div class="form-group{{if .Error}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-lock"></i></span>
                                <input  type="password" class="form-control" name="password" placeholder="Password">
                            </div>
                            <span class="help-block">{{.Error}}</span>
                        </div>
                        {{if .ShowRemember}}
                        <div class="checkbox">
                            <label>
                                <input type="checkbox" name="rm" value="true"> Remember Me
                            </label>
                        </div>
                        {{end}}
                        <button class="btn btn-primary btn-block" type="submit">Login</button>
                        {{if .ShowRecover}}
                        <a class="btn btn-link btn-block" type="submit" href="/recover">Recover Account</a>
                        {{end}}
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>


