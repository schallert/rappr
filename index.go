package main

const INDEX_HTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Rappr</title>

  <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.5.0/pure-min.css">

  <!-- Bad practice, don't care -->
  <style>
    .main {
      padding-left: 50px;
    }
  </style>

</head>
<body>
  <div class="main">
    <h1>Rappr</h1>
    <p>Derek: {{.Derek}}</p>
    <p>Jay: {{.Jay}}</p>
  </div>
</body>
</html>`
