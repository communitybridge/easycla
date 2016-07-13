## Stage

Use the stage plugin to reveal a hidden "shelf" of content, like some common offscreen navigation. Stage will slide in your hidden content—from the left (default) or right—and shift your page contents.

Add `data-toggle="stage"` and a `data-target` to any clickable element that you want to trigger a stage transition. When sliding in from the right, use a negative distance like `data-distance="-250"`.

### Options

- easing: `cubic-bezier(.2,.7,.5,1)`
- duration: 300
- delay: 0
- distance: 250
- hiddenElements: `#sidebar` - hidden elements visibility will be toggle'd on open/close of the stage. This is done for performance reasons in chrome.

### Customizing

Stage shelves—the hidden content areas—can house just about any content, but you'll likely have to tweak some styles depending on what you place within them.

Be sure to match any `data-distance` values with the CSS-based `width` on the `.stage-shelf`.

### Example

For an interactive demo, refer to the `Menu` button at the top of this page. **The example below shows the minimal required markup.**

{% highlight html %}
<div class="stage">
  <button class="btn btn-link stage-toggle" data-target=".stage" data-toggle="stage">
    <span class="icon icon-menu"></span>
    Menu
  </button>

  <div class="stage-shelf">
    <!-- Hidden shelf content -->
  </div>

  <!-- All other page content -->
</div>
{% endhighlight %}

For a right-aligned stage shelf and button, the **minimal required markup looks like this:**

{% highlight html %}
<div class="stage">
  <button class="btn btn-link stage-toggle stage-toggle-right" data-target=".stage" data-toggle="stage" data-distance="-250">
    <span class="icon icon-menu"></span>
    Menu
  </button>

  <div class="stage-shelf stage-shelf-right">
    <!-- Hidden shelf content -->
  </div>

  <!-- All other page content -->
</div>
{% endhighlight %}
