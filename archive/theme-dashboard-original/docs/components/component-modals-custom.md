## Custom Modals

Support for multiple, stackable modal bodies and scrolling bodies using the new `.modal-body-scroller` modifier.

{% example html%}
<div class="modal">
  <div class="modal-dialog modal-md">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title">Medium Modal</h4>
      </div>
      <div class="modal-body modal-body-scroller">
        <p>Maecenas faucibus mollis interdum. Nulla vitae elit libero, a pharetra augue. Aenean eu leo quam. Pellentesque ornare sem lacinia quam venenatis vestibulum. Maecenas faucibus mollis interdum. Donec ullamcorper nulla non metus auctor fringilla. Morbi leo risus, porta ac consectetur ac, vestibulum at eros.</p>

        <p>Nullam id dolor id nibh ultricies vehicula ut id elit. Sed posuere consectetur est at lobortis. Etiam porta sem malesuada magna mollis euismod. Nullam quis risus eget urna mollis ornare vel eu leo.</p>

        <p>Nullam quis risus eget urna mollis ornare vel eu leo. Sed posuere consectetur est at lobortis. Nullam quis risus eget urna mollis ornare vel eu leo. Integer posuere erat a ante venenatis dapibus posuere velit aliquet. Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui.</p>

        <p>Sed posuere consectetur est at lobortis. Maecenas faucibus mollis interdum. Donec sed odio dui. Praesent commodo cursus magna, vel scelerisque nisl consectetur et.</p>

        <p>Nullam id dolor id nibh ultricies vehicula ut id elit. Integer posuere erat a ante venenatis dapibus posuere velit aliquet. Nullam id dolor id nibh ultricies vehicula ut id elit. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.</p>

        <p>Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nullam id dolor id nibh ultricies vehicula ut id elit. Curabitur blandit tempus porttitor. Curabitur blandit tempus porttitor. Maecenas faucibus mollis interdum. Nullam quis risus eget urna mollis ornare vel eu leo.</p>

        <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Nullam quis risus eget urna mollis ornare vel eu leo. Cras justo odio, dapibus ac facilisis in, egestas eget quam. Fusce dapibus, tellus ac cursus commodo, tortor mauris condimentum nibh, ut fermentum massa justo sit amet risus. Nullam quis risus eget urna mollis ornare vel eu leo.</p>

        <p>Nullam quis risus eget urna mollis ornare vel eu leo. Cras mattis consectetur purus sit amet fermentum. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Cras mattis consectetur purus sit amet fermentum.</p>
      </div>

      <div class="modal-body">
        <input type="text" class="form-control" placeholder="Message">
      </div>
    </div>
  </div>
</div>
{% endexample %}

Also, drop the regular footer for a new two-up set of actions using the `.modal-actions` component.

{% example html%}
<div class="modal">
  <div class="modal-dialog modal-sm">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title">Small modal</h4>
      </div>
      <div class="modal-body">
        <p>Modal body text...</p>
      </div>
      <div class="modal-actions">
        <button type="button" class="btn-link modal-action" data-dismiss="modal">Cancel</button>
        <button type="button" class="btn-link modal-action" data-dismiss="modal">
          <strong>Continue</strong>
        </button>
      </div>
    </div>
  </div>
</div>
{% endexample %}
