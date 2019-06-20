// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'member-modal'
})
@Component({
  selector: 'member-modal',
  templateUrl: 'member-modal.html',
})
export class MemberModal {

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService
  ) {
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {

  }

  // MemberModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

}
