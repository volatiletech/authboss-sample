<div class="row">
	<div class="col-md-offset-4 col-md-4">
		<div class="panel panel-default">
			<div class="panel-heading">Reset Password</div>
			<div class="panel-body">
				<form method="POST">
					<div class="form-group {{with .errs}}{{with $errlist := index . "password"}}has-error{{end}}{{end}}">
						<input type="password" class="form-control" name="password" placeholder="Password" value="{{.password}}" />
						{{with .errs}}{{with $errlist := index . "password"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
					</div>
					<div class="form-group {{with .errs}}{{with $errlist := index . "confirm_password"}}has-error{{end}}{{end}}">
						<input type="password" class="form-control" name="confirm_password" placeholder="Confirm Password" value="{{.confirmPassword}}" />
						{{with .errs}}{{with $errlist := index . "confirm_password"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
					</div>
					<input type="hidden" name="token" value="{{.token}}" />
					<input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />
					<div class="row">
						<div class="col-md-offset-1 col-md-10">
							<button class="btn btn-primary btn-block" type="submit">Reset</button>
						</div>
					</div>
					<div class="row">
						<div class="col-md-offset-1 col-md-10">
							<a class="btn btn-link btn-block" href="{{mountpathed "login"}}">Cancel</a>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>