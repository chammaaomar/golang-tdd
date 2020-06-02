package cyoa

import "strings"

var storyTempl = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>GoAdventure!</title>
	</head>
	<body>
		
		<h1>{{.Title}}</h1>
		{{range .Story}}
		<p>{{.}}</p>
		{{end}}

		<ul>
			{{range .Options}}
			<li>
				<a href="%s/{{.Arc}}"><p>{{.Text}}</p>
			</li>
			{{end}}
		</ul>
		</body>
</html>
`

var homeTempl = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>GoAdventure</title>
	</head>
	<body>
		<h1>Stories</h1>
		<ul>
			{{range .}}
			<li><a href="%s/{{.}}/intro"><p>{{. | toTitle}}</p></li>
			{{end}}	
		</ul>
	</body>
</html>
`

// toTitle maps "the-little-blue-gopher" to "The Little Blue Gopher"
func toTitle(in string) string {
	return strings.Title(strings.Replace(in, "-", " ", -1))
}
