<div class="row">
	<div class="col-md-offset-4 col-md-4">
		<div class="panel panel-default">
			<div class="panel-heading">Account Recover</div>
			<div class="panel-body">
				<form method="POST">
					{{$pid := .primaryID}}
					<div class="form-group {{with .errs}}{{with $errlist := index . $pid}}has-error{{end}}{{end}}">
						<input type="text" class="form-control" name="{{.primaryID}}" placeholder="{{title .primaryID}}" value="{{.primaryIDValue}}" />
						{{with .errs}}{{with $errlist := index . $pid}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
					</div>
					{{$cpid := .primaryID | printf "confirm_%s"}}
					<div class="form-group {{with .errs}}{{with $errlist := index . $cpid}}has-error{{end}}{{end}}">
						<input type="text" class="form-control" name="confirm_{{.primaryID}}" placeholder="Confirm {{title .primaryID}}" value="{{.confirmPrimaryIDValue}}" />
						{{with .errs}}{{with $errlist := index . $cpid}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
					</div>
					<input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />
					<div class="row">
						<div class="col-md-offset-1 col-md-10">
							<button class="btn btn-primary btn-block" type="submit">Recover</button>
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