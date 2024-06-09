const mix = require("laravel-mix");

mix
  .setPublicPath("public")
  .postCss("resources/css/input.css", "public/css/main.css", [
    require("tailwindcss"),
  ])
  .js("resources/js/main.js", "public/js/main.js");

if (mix.inProduction()) {
  mix.copy(
    "node_modules/@shoelace-style/shoelace/dist/assets/icons",
    "public/icons"
  );
}
