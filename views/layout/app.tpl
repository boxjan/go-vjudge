<!DOCTYPE html>
<html>
    <head>
        {{ template "layout/component/header.tpl" .}}
    </head>
<body>
<header>
    {{ template "layout/component/navbar.tpl" .}}
</header>
<div id="app">
    <main class="py-4">
        {{ template "main-content" .}}
    </main>
</div>
<footer>
{{ template "layout/footer.tpl" .}}
</footer>
</body>
    {{ template "extant-script" .}}
    <script>
        if (typeof mixin == "undefined") {
            mixin = {};
        }
        let vue = new Vue({
            mixins: [mixin],
            el: '#app',
            data: function() {
                return {
                };
            },
            methods: {
            },
        })
    </script>
</html>