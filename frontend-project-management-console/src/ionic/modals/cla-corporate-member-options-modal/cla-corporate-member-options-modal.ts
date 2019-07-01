// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, ModalController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'cla-corporate-member-options-modal'
})
@Component({
  selector: 'cla-corporate-member-options-modal',
  templateUrl: 'cla-corporate-member-options-modal.html',
})
export class ClaCorporateMemberOptionsModal {
  members: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {
  }

  openClaCorporateWhitelistModal() {
    let modal = this.modalCtrl.create('ClaCorporateWhitelistModal');
    modal.present();
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }



}
