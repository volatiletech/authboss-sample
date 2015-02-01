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
                <div class="panel-heading">Reset Password</div>
                <div class="panel-body">
                    <form action="/recover/complete" method="POST">
                        <input type="hidden" name="token" value="{{.Token}}" />
                        {{$passwordErrs := .ErrMap.password}}
                        <div class="form-group{{if $passwordErrs}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-lock"></i></span>
                                <input class="form-control" type="text" name="password" placeholder="Password" required />
                            </div>
                            {{range $err := $passwordErrs}}
                                <span class="help-block">{{print $err}}</span>
                            {{end}}
                        </div>

                        {{$confirmPasswordErrs := .ErrMap.confirmPassword}}
                        <div class="form-group{{if $confirmPasswordErrs}} has-error{{end}}">
                            <div class="input-group">
                                <span class="input-group-addon"><i class="fa fa-lock"></i></span>
                                <input class="form-control" type="text" name="confirmPassword" placeholder="Confirm Password" required />
                            </div>
                            {{range $err := $confirmPasswordErrs}}
                                <span class="help-block">{{print $err}}</span>
                            {{end}}
                        </div>
                        <button class="btn btn-primary btn-block" type="submit">Submit</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>