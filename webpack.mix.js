const mix = require('laravel-mix');

/*
 |--------------------------------------------------------------------------
 | Mix Asset Management
 |--------------------------------------------------------------------------
 |
 | Mix provides a clean, fluent API for defining some Webpack build steps
 | for your Laravel application. By default, we are compiling the Sass
 | file for the application as well as bundling up all the JS files.
 |
 */

mix.setPublicPath('./');

mix.options({
    processCssUrls: true,
    fileLoaderDirs: {
        fonts: 'static/fonts'
    }
});

mix.js('resources/js/app.js', 'static/js')
    .sass('resources/sass/app.scss', 'static/css');

mix.sass('resources/sass/error.scss', 'static/css');

if (mix.inProduction()) {
    mix.version();
}