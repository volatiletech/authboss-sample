{{define "pagetitle"}}Blogs - Index{{end}}

{{$loggedin := .loggedin}}
{{if $loggedin}}
<div class="row" style="margin-bottom: 20px;">
	<div class="col-md-offset-9 col-md-2 text-right">
		<a class="btn btn-primary" href="/blogs/new"><i class="fa fa-plus"></i> New Post</a>	
	</div>
</div>
{{end}}

<div class="row">
	<div class="col-md-offset-1 col-md-10">
		{{range .posts}}
		<div class="panel panel-info">
			<div class="panel-heading">
				<div class="row">
					<div class="col-md-6">{{.Title}}</div>
					<div class="col-md-6 text-right">
						{{if $loggedin}}
						<a class="btn btn-xs btn-link" href="/blogs/{{.ID}}/edit">Edit</a>
						<a class="btn btn-xs btn-link" href="/blogs/{{.ID}}/destroy">Delete</a>
						{{end}}
					</div>
				</div>
			</div>
			<div class="panel-body">{{.Content}}</div>
			<div class="panel-footer">
				<div class="row">
					<div class="col-md-6">By {{.AuthorID}}</div>
					<div class="col-md-6 text-right">Posted on {{formatDate .Date}}</div>
				</div>
			</div>
		</div>
		{{end}}
	</div>
</div>