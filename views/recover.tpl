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
                    <form action="recover" method="POST">
                        <div class="form-group">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-user"></i></span>
                                <input class="form-control" type="text" name="username" placeholder="Username" value="{{.Username}}" required />
                            </div>
                        </div>
                        <div class="form-group{{if .Error}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-user"></i></span>
                                <input class="form-control" type="text" name="confirmUsername" placeholder="Confirm Username" required />
                            </div>
                            <span class="help-block">{{.Error}}</span>
                        </div>
                        <button class="btn btn-primary btn-block" type="submit">Recover</button>
                        <a class="btn btn-link btn-block" type="submit" href="/login">Cancel</a>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>