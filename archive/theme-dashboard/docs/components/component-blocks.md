## Blocks

Blocks are a brand new composing metaphor exclusively for use with our Marketing toolkit. Build your site up by stacking several blocks on top of each other to provide focused, attention grabbing units of content.

{% example html %}
<div class="block text-center">
  <h1 class="block-title">A basic block</h1>
  <h4 class="text-muted">Use these to package stackable, image driven content.</h4>
  <button class="btn btn-primary m-t">With a simple button</button>
</div>
{% endexample %}

We have several basic block types, including an inverse block.

{% example html %}
<div class="block block-inverse text-center">
  <h1 class="block-title">An inverse block</h1>
  <h4 class="text-muted">Use the inverse modifier for working over dark images.</h4>
  <button class="btn btn-primary m-t">With a simple button</button>
</div>
{% endexample %}

We've also made it easy to integrate embeded content for more interactive block backgrounds.

{% example html %}
<div class="block block-inverse text-center">
  <div class="block-foreground">
    <h1 class="block-title">An embed block</h1>
    <h4 class="text-muted">Use block-background to integrate interactive backgrounds.</h4>
    <button class="btn btn-default btn-outline m-t">With a simple button</button>
  </div>
  <div class="block-background">
    <iframe frameBorder="0" src="https://a.tiles.mapbox.com/v4/jacobthornton.6681fb42/attribution.html?access_token=pk.eyJ1IjoiamFjb2J0aG9ybnRvbiIsImEiOiJlMGRmZmJlNDZkNDhlN2EzMTQ0YWFiNjhlN2RiZWY1ZCJ9.hO-UNIIplnebJYkya-8TEQ"></iframe>
  </div>
</div>
{% endexample %}

Use different modifiers like `block-bordered` and `block-angle` to experiment with different visual treatments and flows between your blocks.

{% example html %}
<div class="block text-center">
  <div class="container-fluid">
    <h4 class="m-b-md">
      Join over 900,000 designers already using Bootstrap. Get Bootstrap <strong>free</strong> forever!
    </h4>
    <form class="form-inline">
      <input class="form-control m-b" placeholder="email address">
      <input class="form-control m-b" type="password" placeholder="Create a Password">
      <button class="btn btn-primary m-b">Get started - free forever</button>
    </form>
    <small class="text-muted">
      By clicking "get started – free Forever!" I agree to Bootstraps
      <a href="#">Terms of service</a>
    </small>
  </div>
</div>
<div class="block block-bordered text-center">
  <div class="container-fluid">
   <blockquote class="pull-quote">
      <p>
        “Notice that simple inset border above. Isn't it lovely.”
      </p>
      <cite>Mark Otto, Huge Nerd</cite>
    </blockquote>
  </div>
</div>
{% endexample %}

Use the `block-fill-height` modifier to make your block fill a user's screen, and then use the responsive alignment classes like `block-xs-middle` or `block-md-bottom` to align your content within the block.

{% example html %}
<div class="block block-fill-height text-center">
  <div class="block-xs-bottom">
    <div class="container-fluid">
     <blockquote class="pull-quote">
        <p>
          “Started at the bottom… etc”
        </p>
        <cite>Drake, OVO</cite>
      </blockquote>
    </div>
  </div>
</div>
{% endexample %}
