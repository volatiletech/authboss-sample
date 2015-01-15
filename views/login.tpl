<html>
<head>
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css" rel="stylesheet">
    <link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<div class="container-fluid">
    <div class="row" style="margin-top: 75px;">
        <div class="col-md-offset-4 col-md-4">
            <div class="panel panel-default">
                <div class="panel-heading">Recover</div>
                <div class="panel-body">
                    <form action="login" method="POST">
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
                        <a class="btn btn-link btn-block" type="submit" href="/recover">Recover Account</a>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>


