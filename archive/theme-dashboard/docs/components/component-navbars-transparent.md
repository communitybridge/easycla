## Transparent navbar

Use the new `.navbar-transparent` variation for placing navbars over backgrounds and interactive content.

{% example html %}
<div class="p-y-lg p-x" style="background: url({{ relative }}assets/img/kanye.jpg) top center; background-size: cover">
  <nav class="navbar navbar-transparent m-b-0">
    <div class="container-fluid">
      <div class="navbar-header">
        <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse-com">
          <span class="sr-only">Toggle navigation</span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
        </button>
        <a class="navbar-brand" href="#">
          <h4 class="text-uppercase m-y-0">Project Name</h4>
        </a>
      </div>
      <div class="navbar-collapse collapse" id="navbar-collapse-com">
        <ul class="nav navbar-nav navbar-right">
          <li class="active"><a href="#">Home</a></li>
          <li><a href="#about">About</a></li>
          <li><a href="#contact">Contact</a></li>
        </ul>
      </div>
    </div>
  </nav>
</div>
{% endexample %}
