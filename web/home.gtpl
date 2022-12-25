<html>
	<head>
	<title></title>
	</head>
	<body>
    <h1>Home</h1>
		Hello {{.Username}} Following: {{.Following}} Followers: {{.Followers}}.
		<form action="/logout" method="post">
			<input type="submit" value="Logout">
		</form>
		<form action="/profile" method="get">
			<input type="submit" value="Profile">
		</form>
		<h3> Follow </h3>
		<form action="/followUser" method="post">
			User's username:<input type="text" name="username">
			<input type="submit" value="Follow User">
		</form>
		<h3> View Other Profile </h3>
		<form action="/otherUser" method="GET">
			User's username:<input type="text" name="id">
			<input type="submit" value="View User">
		</form>
		<div style="width:100%; height:10%">
		<h3> Post </h3>
		<form action="/createPost" method="post">
			Post Content:<input type="text" name="content">
			<input type="submit" value="Create Post">
		</form>
		<div style="width:100%; height:10%">
		<h3> Feed </h3>
		{{if .Posts}}
			{{range .Posts}}
				<div style="border: thin solid black">
				<h3 style="display: inline-block;">{{.author}}</h3> <p style="display: inline-block;">{{.createdAt}}</p>
				<p>{{.content}}</p>
				</div>
			{{end}}
		{{else}}
			<p>Your feed is empty</p>
		{{end}}
	</body>
</html>