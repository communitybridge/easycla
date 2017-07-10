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
  orgName: string;
  memberships: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.orgName = navParams.get('orgName');
    this.memberships = navParams.get('memberships');
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
