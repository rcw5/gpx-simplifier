<html>
<head>
       <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="http://127.0.0.1:8081/upload" method="post">
    GPX File to upload:
    <input type="file" name="uploadfile" />
    <br/>
    Number of files:
    <select name="num_files" id="num_files">
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
    <br />
    Points per file:
    <select name="points_per_file" id="points_per_file">
    <option value="250">250</option>
    <option value="500">500</option>
    <option value="750">750</option>
    <option value="1000">1000</option>
    </select>
    <br />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
</body>
</html>
