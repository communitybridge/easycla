## List group

Bootstrap's default list group component is extended with a few additional features.

### List group header

New to the theme is a header for list groups, similar to panels. use it for standalone lists as needed. This should not be used with panels.

{% example html %}
<ul class="list-group">
  <li class="list-group-header">List group header</li>
  <li class="list-group-item">Cras justo odio</li>
  <li class="list-group-item">Dapibus ac facilisis in</li>
  <li class="list-group-item">Morbi leo risus</li>
  <li class="list-group-item">Porta ac consectetur ac</li>
  <li class="list-group-item">Vestibulum at eros</li>
</ul>
{% endexample %}


### List group progress

Similar to stat lists, add a background progress indicator to any list group item.

{% example html %}
<ul class="list-group">
  <li class="list-group-item">
    Cras justo odio
    <span class="list-group-progress" style="width: 75%;"></span>
  </li>
  <li class="list-group-item">
    Porta ac consectetur ac
    <span class="list-group-progress" style="width: 50%;"></span>
  </li>
  <li class="list-group-item">
    Vestibulum at eros
    <span class="list-group-progress" style="width: 25%;"></span>
  </li>
</ul>
{% endexample %}
