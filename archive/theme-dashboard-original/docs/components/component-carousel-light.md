## Carousel Light

The light carousel is a modified version of Bootstrap's default carousel, using lighter styles on all carousel controls.

{% example html %}
<div id="carousel-example-generic-2" class="carousel carousel-light slide" data-ride="carousel">
  <ol class="carousel-indicators">
    <li data-target="#carousel-example-generic-2" data-slide-to="0" class="active"></li>
    <li data-target="#carousel-example-generic-2" data-slide-to="1"></li>
    <li data-target="#carousel-example-generic-2" data-slide-to="2"></li>
  </ol>
  <div class="carousel-inner" role="listbox">
    <div class="item active">
      <img src="http://placehold.it/1140x500/fff/333" alt="First slide">
    </div>
    <div class="item">
      <img src="http://placehold.it/1140x500/fff/333" alt="Second slide">
    </div>
    <div class="item">
      <img src="http://placehold.it/1140x500/fff/333" alt="Third slide">
    </div>
  </div>
  <a class="left carousel-control" href="#carousel-example-generic-2" role="button" data-slide="prev">
    <span class="icon icon-chevron-thin-left" aria-hidden="true"></span>
    <span class="sr-only">Previous</span>
  </a>
  <a class="right carousel-control" href="#carousel-example-generic-2" role="button" data-slide="next">
    <span class="icon icon-chevron-thin-right" aria-hidden="true"></span>
    <span class="sr-only">Next</span>
  </a>
</div>
{% endexample %}
