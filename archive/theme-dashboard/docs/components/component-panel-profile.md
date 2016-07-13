## Panel profile

Feature a user's profile with the `.panel-profile` variant.

{% example html %}
<div class="m-t">
  <div class="row">
    <div class="col-md-6">
      <div class="panel panel-default panel-profile">
        <div class="panel-heading" style="background-image: url(https://igcdn-photos-h-a.akamaihd.net/hphotos-ak-xfa1/t51.2885-15/11312291_348657648678007_1202941362_n.jpg);"></div>
        <div class="panel-body text-center">
          <img class="panel-profile-img" src="{{ relative }}assets/img/avatar-fat.jpg">
          <h5 class="panel-title">Jacob Thornton</h5>
          <p class="m-b-md">Big belly rude boy, million dollar hustler. Unemployed.</p>
          <button class="btn btn-primary-outline btn-sm">
            <span class="icon icon-add-user"></span> Follow
          </button>
        </div>
      </div>
    </div>
    <div class="col-md-6">
      <div class="panel panel-default panel-profile">
        <div class="panel-heading" style="background-image: url(https://igcdn-photos-h-a.akamaihd.net/hphotos-ak-xaf1/t51.2885-15/11240760_365538423656311_112029877_n.jpg);"></div>
        <div class="panel-body text-center">
          <img class="panel-profile-img" src="{{ relative }}assets/img/avatar-mdo.png">
          <h5 class="panel-title">Mark Otto</h5>
          <p class="m-b-md">GitHub and Bootstrap. Formerly at Twitter. Huge nerd.</p>
          <button class="btn btn-primary-outline btn-sm">
            <span class="icon icon-add-user"></span> Follow
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
{% endexample %}
