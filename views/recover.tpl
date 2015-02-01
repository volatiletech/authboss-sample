<html>
<head>
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css" rel="stylesheet">
    <link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js" type="text/javascript"></script>
</head>
<body>
<div class="container-fluid">
    {{if .FlashError}}
    <div class="row">
        <div class="col-xs-offset-3 col-md-6">
            <div class="alert alert-danger alert-dismissable" style="margin-top: 75px;">
                <button type="button" class="close" data-dismiss="alert"><span>&times;</span></button>
                {{print .FlashError}}
            </div>
        </div>
    </div>
    {{end}}
    <div class="row" style="margin-top: 75px;">
        <div class="col-md-offset-4 col-md-4">
            <div class="panel panel-default">
                <div class="panel-heading">Recover Account</div>
                <div class="panel-body">
                    <form action="/recover" method="POST">
                        {{$usernameErrs := .ErrMap.username}}
                        <div class="form-group{{if $usernameErrs}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-user"></i></span>
                                <input class="form-control" type="text" name="username" placeholder="Username" value="{{.Username}}" />
                            </div>
                            {{range $err := $usernameErrs}}
                                <span class="help-block">{{print $err}}</span>
                            {{end}}
                        </div>

                        {{$confirmUsernameErrs := .ErrMap.confirmUsername}}
                        <div class="form-group{{if $confirmUsernameErrs}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-user"></i></span>
                                <input class="form-control" type="text" name="confirmUsername" placeholder="Confirm Username" value="{{.ConfirmUsername}}" />
                            </div>
                            {{range $err := $confirmUsernameErrs}}
                                <span class="help-block">{{print $err}}</span>
                            {{end}}
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