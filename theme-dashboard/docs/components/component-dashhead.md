## Dashhead

The dashhead is a custom component built to house all the textual headings, form controls, buttons, and more that are common for the top of dashboard page.

{% example html %}
<div class="dashhead">
  <div class="dashhead-titles">
    <h6 class="dashhead-subtitle">Dashboards</h6>
    <h3 class="dashhead-title">Overview</h3>
  </div>

  <div class="dashhead-toolbar">
    <div class="input-with-icon dashhead-toolbar-item">
      <input type="text" value="01/01/15 - 01/08/15" class="form-control" data-provide="datepicker">
      <span class="icon icon-calendar"></span>
    </div>
    <span class="dashhead-toolbar-divider hidden-xs"></span>
    <div class="btn-group dashhead-toolbar-item btn-group-thirds">
      <button type="button" class="btn btn-primary-outline">Day</button>
      <button type="button" class="btn btn-primary-outline active">Week</button>
      <button type="button" class="btn btn-primary-outline">Month</button>
    </div>
  </div>
</div>
{% endexample %}
