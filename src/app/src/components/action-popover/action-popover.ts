import { Component, Output, EventEmitter } from '@angular/core';
import { NavParams, ViewController } from 'ionic-angular';

@Component({
  template: `
    <ion-list class="action-popover">
      <button ion-item *ngFor="let item of popoverItems; let index = index;" (click)='popoverAction(index)'>
        {{ item.label }}
      </button>
    </ion-list>
  `
})
export class ActionPopover {

  popoverItems: any;

  @Output() popoverNotice: EventEmitter<{}> = new EventEmitter<{}>();

  constructor(private navParams: NavParams, private viewCtrl: ViewController) {

  }

  ngOnInit() {
    if (this.navParams.data) {
      this.popoverItems = this.navParams.data.items
    }
  }

  popoverAction(index) {
    let callback = this.popoverItems[index].callback;
    let callbackData = this.popoverItems[index].callbackData;
    this.viewCtrl.dismiss({
      callback: callback,
      callbackData: callbackData
    });
  }

}
