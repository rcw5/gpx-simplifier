<html>
<head>
  <!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<!-- Optional theme -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

       <title>The Wonderful and Fantastic GPX File Splitter and Simplifier</title>
</head>
<body>

    <div class="container">
        <h1>GPX Simplifier</h1>
        {{if .Success}}
        <p>Your file is ready to download: <a download="{{.FileName}}.zip" href="{{.URL}}">here</a>.</p>
        {{else}}
        <p class="bg-danger" style="padding: 10px"><b>Process failed</b>: {{ .ErrorMessage }}</p>
        <p><a href="/upload">Try again?</a></p>
        {{end}}
        <h6>Designed and developed by <a href="http://www.github.com/rcw5/">Robin Watkins</a></h6>
    </div>


</body>
</html>
