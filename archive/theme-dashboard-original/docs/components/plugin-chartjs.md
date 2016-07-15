## Chart.js

Charts and graphs are available through [Chart.js](http://www.chartjs.org), a clean and responsive canvas-based chart rendering JavaScript library. For even easier integration, we've created a small plugin for use with this toolkit for rendering four types of graphs: doughnut, line, bar, and the custom sparkline.

We recommend reviewing the full [Chart.js documentation](http://www.chartjs.org/docs/) as you implement or modify any charts here.

### Data API

The data API allows you to use Chart.js arguments and options by writing HTML. Take any option provided by Chart.js and simply hypenate it as an HTML attribute. For example, the doughnut's `segmentStrokeColor` option, becomes `data-segment-stroke-color`.

Below are examples of the graph types we support out of the box in this toolkit, implemented with the data API approach.

### Doughnut

{% example html %}
<div class="w-sm m-x-auto">
  <canvas
    class="ex-graph"
    width="200" height="200"
    data-chart="doughnut"
    data-value="[{ value: 230, color: '#1ca8dd', label: 'Returning' }, { value: 130, color: '#1bc98e', label: 'New' }]"
    data-segment-stroke-color="#252830">
  </canvas>
</div>
{% endexample %}

### Bar

{% example html %}
<div>
  <canvas
    class="ex-line-graph"
    width="600" height="400"
    data-chart="bar"
    data-scale-line-color="transparent"
    data-scale-grid-line-color="rgba(255,255,255,.05)"
    data-scale-font-color="#a2a2a2"
    data-labels="['August','September','October','November','December','January','February']"
    data-value="[{ label: 'First dataset', data: [65, 59, 80, 81, 56, 55, 40] }, { label: 'Second dataset', data: [28, 48, 40, 19, 86, 27, 90] }]">
  </canvas>
</div>
{% endexample %}

### Line

{% example html %}
<div>
  <canvas
    class="ex-line-graph"
    data-chart="line"
    data-scale-line-color="transparent"
    data-scale-grid-line-color="rgba(255,255,255,.05)"
    data-scale-font-color="#a2a2a2"
    data-labels="['','Aug 29','','','Sept 5','','','Sept 12','','','Sept 19','']"
    data-value="[{fillColor: 'rgba(28,168,221,.03)', data: [2500, 3300, 2512, 2775, 2498, 3512, 2925, 4275, 3507, 3825, 3445, 3985]}]">
  </canvas>
</div>
{% endexample %}

### Sparkline

{% example html %}
<div class="row">
  <div class="col-sm-6 col-md-4">
    <div class="statcard statcard-success">
      <div class="p-a">
        <span class="statcard-desc">Page views</span>
        <h2 class="statcard-number">
          12,938
          <small class="delta-indicator delta-positive">5%</small>
        </h2>
        <hr class="statcard-hr m-b-0">
      </div>
      <canvas
        class="sparkline"
        data-chart="spark-line"
        data-value="[{data:[28,68,41,43,96,45,100]}]"
        data-labels="['a','b','c','d','e','f','g']"
        width="378" height="94">
      </canvas>
    </div>
  </div>
</div>
{% endexample %}
