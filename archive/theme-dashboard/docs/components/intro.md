# Intro

Hey there! You're looking at the docs for an Official Bootstrap Theme—thanks for your purchase! This theme has been lovingly crafted by Bootstrap's founders and friends to help you build awesome projects on the web. Let's dive in.

Each theme is designed as it's own toolkit—a series of well designed, intuitive, and cohesive components—built on top of Bootstrap. If you've used Bootstrap itself, you'll find this all super familiar, but with new aesthetics, new components, beautiful and extensive examples, and easy-to-use build tools and documentation.

Since this theme is based on Bootstrap, and includes nearly everything Bootstrap itself does, you'll want to give the [official project documentation](http://getbootstrap.com) a once over before continuing. There's also [the kitchen sink]({{ relative }}bootstrap/index.html)—a one-page view of all Bootstrap's components restyled by this theme.

For everything else, including compiling and using this Bootstrap theme, read through the docs below.

Thanks, and enjoy!

# What's included

Within your Bootstrap theme you'll find the following directories and files, grouping common resources and providing both compiled and minified distribution files, as well as raw source files.

{% highlight bash %}
theme/
  ├── gulpfile.js
  ├── package.json
  ├── README.md
  ├── docs/
  ├── less/
  │   ├── bootstrap/
  │   ├── custom/
  │   ├── variables.less
  │   └── toolkit.less
  ├── js/
  │   ├── bootstrap/
  │   └── custom/
  ├── fonts/
  │   ├── bootstrap-entypo.eot
  │   ├── bootstrap-entypo.svg
  │   ├── bootstrap-entypo.ttf
  │   ├── bootstrap-entypo.woff
  │   └── bootstrap-entypo.woff2
  └── dist/
      ├── toolkit.css
      ├── toolkit.css.map
      ├── toolkit.min.css
      ├── toolkit.min.css.map
      ├── toolkit.js
      └── toolkit.min.js
{% endhighlight %}


## Getting started

To view your Bootstrap Theme documentation, simply find the docs directory and open index.html in your favorite browser.

{% highlight bash %}
$ open docs/index.html
{% endhighlight %}


## Gulpfile.js

If you're after more customization we've also included a custom [Gulp](http://gulpjs.com) file, which can be used to quickly re-compile a theme's CSS and JS. You'll need to install both [Node](https://nodejs.org/en/download/) and Gulp before using our included `gulpfile.js`.

Once node is installed, run the following npm command to install Gulp.

{% highlight bash %}
$ npm install gulp -g
{% endhighlight %}

When you're done, make sure you've installed the rest of the theme's dependencies:

{% highlight bash %}
$ npm install
{% endhighlight %}

Now, modify your source files and run `gulp` to generate new local `dist/` files automatically. **Be aware that this replaces existing `dist/` files**, so proceed with caution.

## Theme source code

The `less/`, `js/`, and `fonts/` directories contain the source code for our CSS, JS, and icon fonts (respectively). Within the `less/` and `js/` directories you'll find two subdirectories:

- `bootstrap/`, which contains the most recently released version of Bootstrap (v3.3.5).
- `custom/`, which contains all of the custom components and overrides authored specifically for this theme.

The `dist/` folder includes everything above, built into single CSS and JS files that can easily be integrated into your project.

The `docs/` folder includes the source code for our documentation, as well as a handful of live examples.

The remaining files not specifically mentioned above provide support for packages, license information, and development.


# Custom builds

Leverage the included source files and `gulpfile.js` to customize your Bootstrap Theme for your exact needs. Change variables, exclude components, and more.

- `toolkit.less` is the entry point for Less files - to build your own custom build, simply modify your local custom files or edit the includes listed here. Note: some themes also rely on a shared `components.less` file, which you can find imported in your `toolkit.less`.
- `variables.less` is home to your theme's variables. Note that your theme's `variables` file depends on and overrides an existing Bootstrap variable file (found in /less/bootstrap/variables.less).


# Basic template

The basic template is a guideline for how to structure your pages when building with a Bootstrap Theme. Included below are all the necessary bits for using the theme's CSS and JS, as well as some friendly reminders.

Copy the example below into a new HTML file to get started with it.

{% highlight html %}
<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- These meta tags come first. -->
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Bootstrap Theme Example</title>

    <!-- Include the CSS -->
    <link href="dist/toolkit.min.css" rel="stylesheet">

  </head>
  <body>
    <h1>Hello, world!</h1>

    <!-- Include jQuery (required) and the JS -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
    <script src="dist/toolkit.min.js"></script>
  </body>
</html>
{% endhighlight %}

# Utilities

Utilities, or utility classes, are a series of single-purpose, immutable classes designed to reduce duplication in your compiled CSS. Each utility takes a common CSS property-value pair and turns it into a reusable class.

These utilities are in addition to those [provided in Bootstrap](http://getbootstrap.com/css/).

## Positioning

{% highlight scss %}
.pos-r { position: relative !important; }
.pos-a { position: absolute !important; }
.pos-f { position: fixed !important; }
{% endhighlight %}

## Width

{% highlight scss %}
.w-sm   { width: 25% !important; }
.w-md   { width: 50% !important; }
.w-lg   { width: 75% !important; }
.w-full { width: 100% !important; }
{% endhighlight %}

## Margin and padding

Assign `margin` or `padding` to an element or a subset of it's sides with shorthand classes. Includes support for individual properties, all properties, and vertical and horizontal properties. All classes are multiples on the global default value, `20px`.

{% highlight scss %}
.m-a-0 { margin:        0 !important; }
.m-t-0 { margin-top:    0 !important; }
.m-r-0 { margin-right:  0 !important; }
.m-b-0 { margin-bottom: 0 !important; }
.m-l-0 { margin-left:   0 !important; }

.m-a { margin:        @spacer !important; }
.m-t { margin-top:    @spacer-y !important; }
.m-r { margin-right:  @spacer-x !important; }
.m-b { margin-bottom: @spacer-y !important; }
.m-l { margin-left:   @spacer-x !important; }
.m-x { margin-right:  @spacer-x !important; margin-left: @spacer-x !important; }
.m-y { margin-top:    @spacer-y !important; margin-bottom: @spacer-y !important; }
.m-x-auto { margin-right: auto !important; margin-left: auto !important; }

.m-a-md { margin:        (@spacer-y * 1.5) !important; }
.m-t-md { margin-top:    (@spacer-y * 1.5) !important; }
.m-r-md { margin-right:  (@spacer-y * 1.5) !important; }
.m-b-md { margin-bottom: (@spacer-y * 1.5) !important; }
.m-l-md { margin-left:   (@spacer-y * 1.5) !important; }
.m-x-md { margin-right:  (@spacer-x * 1.5) !important; margin-left:   (@spacer-x * 1.5) !important; }
.m-y-md { margin-top:    (@spacer-y * 1.5) !important; margin-bottom: (@spacer-y * 1.5) !important; }

.m-a-lg { margin:        (@spacer-y * 3) !important; }
.m-t-lg { margin-top:    (@spacer-y * 3) !important; }
.m-r-lg { margin-right:  (@spacer-y * 3) !important; }
.m-b-lg { margin-bottom: (@spacer-y * 3) !important; }
.m-l-lg { margin-left:   (@spacer-y * 3) !important; }
.m-x-lg { margin-right:  (@spacer-x * 3) !important; margin-left:   (@spacer-x * 3) !important; }
.m-y-lg { margin-top:    (@spacer-y * 3) !important; margin-bottom: (@spacer-y * 3) !important; }
{% endhighlight %}

{% highlight scss %}
.p-a-0 { padding:        0 !important; }
.p-t-0 { padding-top:    0 !important; }
.p-r-0 { padding-right:  0 !important; }
.p-b-0 { padding-bottom: 0 !important; }
.p-l-0 { padding-left:   0 !important; }

.p-a { padding:        @spacer !important; }
.p-t { padding-top:    @spacer-y !important; }
.p-r { padding-right:  @spacer-x !important; }
.p-b { padding-bottom: @spacer-y !important; }
.p-l { padding-left:   @spacer-x !important; }
.p-x { padding-right:  @spacer-x !important; padding-left:   @spacer-x !important; }
.p-y { padding-top:    @spacer-y !important; padding-bottom: @spacer-y !important; }

.p-a-md { padding:        (@spacer-y * 1.5) !important; }
.p-t-md { padding-top:    (@spacer-y * 1.5) !important; }
.p-r-md { padding-right:  (@spacer-y * 1.5) !important; }
.p-b-md { padding-bottom: (@spacer-y * 1.5) !important; }
.p-l-md { padding-left:   (@spacer-y * 1.5) !important; }
.p-x-md { padding-right:  (@spacer-x * 1.5) !important; padding-left:   (@spacer-x * 1.5) !important; }
.p-y-md { padding-top:    (@spacer-y * 1.5) !important; padding-bottom: (@spacer-y * 1.5) !important; }

.p-a-lg { padding:        (@spacer-y * 3) !important; }
.p-t-lg { padding-top:    (@spacer-y * 3) !important; }
.p-r-lg { padding-right:  (@spacer-y * 3) !important; }
.p-b-lg { padding-bottom: (@spacer-y * 3) !important; }
.p-l-lg { padding-left:   (@spacer-y * 3) !important; }
.p-x-lg { padding-right:  (@spacer-x * 3) !important; padding-left:   (@spacer-x * 3) !important; }
.p-y-lg { padding-top:    (@spacer-y * 3) !important; padding-bottom: (@spacer-y * 3) !important; }
{% endhighlight %}

## Responsive text alignment

Use these classes to easily switch `text-align` across viewports when designing responsive pages.

{% highlight scss %}
.text-xs-left   { text-align: left; }
.text-xs-right  { text-align: right; }
.text-xs-center { text-align: center; }

@media (min-width: @screen-sm-min) {
  .text-sm-left   { text-align: left; }
  .text-sm-right  { text-align: right; }
  .text-sm-center { text-align: center; }
}

@media (min-width: @screen-md-min) {
  .text-md-left   { text-align: left; }
  .text-md-right  { text-align: right; }
  .text-md-center { text-align: center; }
}

@media (min-width: @screen-lg-min) {
  .text-lg-left   { text-align: left; }
  .text-lg-right  { text-align: right; }
  .text-lg-center { text-align: center; }
}
{% endhighlight %}
