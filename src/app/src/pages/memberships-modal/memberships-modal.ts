import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, IonicPage } from 'ionic-angular';

@IonicPage({
  segment: 'memberships-modal'
})
@Component({
  selector: 'memberships-modal',
  templateUrl: 'memberships-modal.html',
})
export class MembershipsModal {
  memberships: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.memberships = navParams.get('memberships');
    console.log(this.memberships);
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {

  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      projectId: projectId,
    });
  }

}
