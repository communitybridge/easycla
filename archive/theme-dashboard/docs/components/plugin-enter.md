## Enter

Enter is a brand new plugin built to transition elements into view on scroll and offer subtle visual flourishes. Simply add a `data-transition="entrance"` attribute and a `transform` style to any element that you'd like to *enter* in when a user scrolls the element into view.

### Options

- easing: `cubic-bezier(.2,.7,.5,1)`,
- duration: 1200,
- delay: 0

### JavaScript API

{% highlight js %}
$('.js-enter').enter()
{% endhighlight %}

### Data API

{% example html %}
<div style="overflow: hidden">
  <div data-transition="entrance" style="transform: translateY(50px)">
    <p>
      Etiam porta sem malesuada magna mollis euismod. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce dapibus, tellus ac cursus commodo, tortor mauris condimentum nibh, ut fermentum massa justo sit amet risus.
    </p>
    <p>
      Donec sed odio dui. Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Curabitur blandit tempus porttitor. Aenean eu leo quam. Pellentesque ornare sem lacinia quam venenatis vestibulum.
    </p>
  </div>
</div>
{% endexample %}
