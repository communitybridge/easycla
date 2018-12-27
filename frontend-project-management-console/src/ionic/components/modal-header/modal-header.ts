import { Component } from "@angular/core";
import { Events } from 'ionic-angular';

@Component({
  selector: "modal-header",
  templateUrl: "modal-header.html"
})
export class ModalHeaderComponent {

  constructor (public events: Events) {}

  triggerDismissEvent () {
    this.events.publish('modal:close');
  }
}
