<div class="page-header p-t-md" id="spacers">
  <h1>Spacer utilities</h1>
</div>

<p>Assign <code>margin</code> or <code>padding</code> to an element or a subset of it's sides with shorthand classes. Includes support for individual properties, all properties, and vertical and horizontal properties. All classes are multiples on the global default value, <code>20px</code>.</p>

<p><strong>All properties include <code>!important</code> in the source files.</strong></p>

<h3>Position</h3>

{% highlight scss %}
.pos-r { position: relative !important; }
.pos-a { position: absolute !important; }
.pos-f { position: fixed !important; }
{% endhighlight %}

<h3>Width</h3>

{% highlight scss %}
.w-sm   { width: 25% !important; }
.w-md   { width: 50% !important; }
.w-lg   { width: 75% !important; }
.w-full { width: 100% !important; }
{% endhighlight %}

<h3>Margin</h3>

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

<h3>Padding</h3>

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
