## Tablesorter

Including in this theme is [Tablesorter](http://tablesorter.com/), a jQuery plugin for easy column sorting on tables. Basic styles for the directional arrows when sorting are included here for ease. Consult the Tablesorter docs for usage and additional customizations.

<div class="docs-example">
  <div class="table-responsive">
    <table class="table" data-sort="table">
      <thead>
        <tr>
          <th><input type="checkbox" class="select-all" id="selectAll"></th>
          <th>Order</th>
          <th>Customer name</th>
          <th>Description</th>
          <th>Date</th>
          <th>Total</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><input type="checkbox" class="select-row"></td>
          <td><a href="#">#10001</a></td>
          <td>First Last</td>
          <td>Admin theme, marketing theme</td>
          <td>01/01/2015</td>
          <td>$200.00</td>
        </tr>
        <tr>
          <td><input type="checkbox" class="select-row"></td>
          <td><a href="#">#10002</a></td>
          <td>Firstname Lastname</td>
          <td>Admin theme</td>
          <td>01/01/2015</td>
          <td>$100.00</td>
        </tr>
        <tr>
          <td><input type="checkbox" class="select-row"></td>
          <td><a href="#">#10003</a></td>
          <td>Name Another</td>
          <td>Personal blog theme</td>
          <td>01/01/2015</td>
          <td>$100.00</td>
        </tr>
        <tr>
          <td><input type="checkbox" class="select-row"></td>
          <td><a href="#">#10004</a></td>
          <td>One More</td>
          <td>Marketing theme, personal blog theme, admin theme</td>
          <td>01/01/2015</td>
          <td>$300.00</td>
        </tr>
        <tr>
          <td><input type="checkbox" class="select-row"></td>
          <td><a href="#">#10005</a></td>
          <td>Name Right Here</td>
          <td>Personal blog theme, admin theme</td>
          <td>01/02/2015</td>
          <td>$200.00</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>

Enabling Tablesorter is easy with jQuery:

{% highlight js %}
$(document).ready(function() {
  // call the tablesorter plugin
  $("[data-sort=table]").tablesorter({
    // Sort on the second column, in ascending order
    sortList: [[1,0]]
  });
});
{% endhighlight %}
