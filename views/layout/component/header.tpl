<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1, user-scalable=no">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="theme-color" content="#FFFFFF">

<!-- CSRF Token -->
<meta name="csrf-token" content="{{ .xsrf_token }}">

<title>{{ .title  }}</title>

<script src="{{ `js/app.js` | asset }}" ></script>
<link href="{{ `css/app.css` | asset }}" rel="stylesheet">
{{ template "extant-header" .}}
