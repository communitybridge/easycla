import { Component } from '@angular/core';

import { NavController, NavParams, ViewController } from 'ionic-angular';

@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html'
})
export class ContactUpdate {
  project: any;
  member: any;
  contact: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController
  ) {
    this.project=this.navParams.get('project');
    console.log(this.project);
    this.member=this.navParams.get('member');
    console.log(this.member);
    this.contact=this.navParams.get('contact');
    console.log(this.contact);
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
