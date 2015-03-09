<div class="row">
	<div class="col-md-offset-4 col-md-4">
		<div class="panel panel-default">
			<div class="panel-heading">Please sign in to continue!</div>
			<div class="panel-body">
				{{if .error}}
				<div class="alert alert-danger">{{.error}}</div>
				{{end}}
				<form method="POST">
					<div class="form-group">
						<input type="text" class="form-control" name="{{.primaryID}}" placeholder="{{title .primaryID}}" value="{{.primaryIDValue}}">
					</div>
					<div class="form-group">
						<input  type="password" class="form-control" name="password" placeholder="Password">
					</div>
					{{if .showRemember}}
					<div class="form-group">
						<input type="checkbox" name="rm" value="true"> Remember Me
					</div>
					{{end}}
					<input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />
					<div class="row">
						<div class="col-md-offset-1 col-md-10">
							<button class="btn btn-primary btn-block" type="submit">Login</button>
						</div>
					</div>
					{{if .showRecover}}
					<div class="row">
						<div class="col-md-offset-1 col-md-10">
							<a class="btn btn-link btn-block" href="{{mountpathed "recover"}}">Recover Account</a>
						</div>
					</div>
					{{end}}
				</form>
			</div>
		</div>
	</div>
</div>
