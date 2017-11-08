import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { PopoverController } from 'ionic-angular';

@IonicPage({
  segment: 'cla-organization-app-modal'
})
@Component({
  selector: 'cla-organization-app-modal',
  templateUrl: 'cla-organization-app-modal.html',
  providers: [CincoService]
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

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
