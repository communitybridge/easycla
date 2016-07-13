## Media list users

{% example html %}
<ul class="media-list media-list-users list-group">
  <li class="list-group-item">
    <div class="media">
      <a class="media-left" href="#">
        <img class="media-object img-circle" src="{{ relative }}assets/img/avatar-fat.jpg">
      </a>
      <div class="media-body">
        <button class="btn btn-primary-outline btn-sm pull-right">
          <span class="icon icon-add-user"></span> Follow
        </button>
        <strong>Jacob Thornton</strong>
        <small>@fat - San Francisco</small>
      </div>
    </div>
  </li>
  <li class="list-group-item">
    <div class="media">
      <a class="media-left" href="#">
        <img class="media-object img-circle" src="{{ relative }}assets/img/avatar-dhg.png">
      </a>
      <div class="media-body">
        <button class="btn btn-primary-outline btn-sm pull-right">
          <span class="icon icon-add-user"></span> Follow
        </button>
        <strong>Dave Gamache</strong>
        <small>@dhg - Palo Alto</small>
      </div>
    </div>
  </li>
  <li class="list-group-item">
    <div class="media">
      <a class="media-left" href="#">
        <img class="media-object img-circle" src="{{ relative }}assets/img/avatar-mdo.png">
      </a>
      <div class="media-body">
        <button class="btn btn-primary-outline btn-sm pull-right">
          <span class="icon icon-add-user"></span> Follow
        </button>
        <strong>Mark Otto</strong>
        <small>@mdo - San Francisco</small>
      </div>
    </div>
  </li>
</ul>
{% endexample %}
