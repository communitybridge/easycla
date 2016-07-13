## Nav bordered

The bordered nav builds on Bootstrap's `.nav` base styles with a new, bolder variation to nav links.

{% example html %}
<ul class="nav nav-bordered">
  <li class="active">
    <a href="#">Bold</a>
  </li>
  <li>
    <a href="#">Minimal</a>
  </li>
  <li>
    <a href="#">Components</a>
  </li>
  <li>
    <a href="#">Docs</a>
  </li>
</ul>
{% endexample %}

Bordered nav links can also be stacked:

{% example html %}
<ul class="nav nav-bordered nav-stacked">
  <li class="nav-header">Examples</li>
  <li class="active">
    <a href="#">Bold</a>
  </li>
  <li><a href="#">Minimal</a></li>

  <li class="nav-divider"></li>
  <li class="nav-header">Sections</li>

  <li><a href="#">Grid</a></li>
  <li><a href="#">Pricing</a></li>
  <li><a href="#">About</a></li>
</ul>
{% endexample %}
