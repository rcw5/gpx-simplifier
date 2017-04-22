<html>
<head>
  <!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<!-- Optional theme -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<!-- Latest compiled and minified JavaScript -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
       <title>The Wonderful and Fantastic GPX File Splitter and Simplifier</title>
</head>
<body>

    <div class="container">
        <h1>GPX Simplifier</h1>
        <p>Use this site to split and simplify a GPX file. The input file will be split into <i>n</i> files, each with a maximum of <i>m</i> trackpoints.
        <form enctype="multipart/form-data" action="http://127.0.0.1:8081/upload" method="post">
          <div class="form-group">
            <label for="file">GPX File to upload:</label>
            <input type="file" name="uploadfile" class="form-control" />
            <p class="help-block">Note: File must be a valid GPX <b>track</b> or it will be rejected.</p>
          </div>
          <div class="form-group">
            <label for="num_files">Number of files:</label>
            <select name="num_files" id="num_files" class="form-control">
            <option value="1">1</option>
            <option value="2">2</option>
            <option value="3">3</option>
            <option value="4">4</option>
            <option value="5">5</option>
            <option value="6">6</option>
            <option value="7">7</option>
            <option value="8">8</option>
            <option value="9">9</option>
            <option value="10">10</option>
            </select>
          </div>
          <div class="form-group">
            <label for="file">Points per file:</label>
            <select name="points_per_file" id="points_per_file" class="form-control">
            <option value="10">10</option>
            <option value="250">250</option>
            <option value="500" selected>500</option>
            <option value="750">750</option>
            <option value="1000">1000</option>
            </select>
          </div>
            <input type="hidden" name="token" value="{{.}}"/>
            <input type="submit" value="Simplify" class="btn btn-default" />
        </form>
        <h6>Designed and developed by <a href="http://www.github.com/rcw5/">Robin Watkins</a></h6>
    </div>


</body>
</html>
