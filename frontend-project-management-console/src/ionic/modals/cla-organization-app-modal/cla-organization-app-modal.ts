import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { PopoverController } from 'ionic-angular';

@IonicPage({
  segment: 'cla-organization-app-modal'
})
@Component({
  selector: 'cla-organization-app-modal',
  templateUrl: 'cla-organization-app-modal.html',
})
export class ClaOrganizationAppModal {

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {

  }

  // TODO: Do we want a call to cla that polls for the installation status?

  openAppPage() {
    window.open('https://github.com/apps/contributor-license-agreement', '_blank');
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
