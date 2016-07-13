## Flex table

Use the flex table classes when you need to keep related elements on the same row, but with flexible individual widths.

{% example html %}
<div class="flextable">
  <div class="flextable-item flextable-primary">
    <input type="text" class="form-control" placeholder="Search orders">
  </div>
  <div class="flextable-item">
    <div class="btn-group">
      <button type="button" class="btn btn-primary-outline">
        <span class="icon icon-pencil"></span>
      </button>
      <button type="button" class="btn btn-primary-outline">
        <span class="icon icon-erase"></span>
      </button>
    </div>
  </div>
</div>
{% endexample %}
