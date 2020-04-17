// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, ViewController, IonicPage, Events, AlertController } from 'ionic-angular';
import { PopoverController } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';
import { ClaService } from '../../services/cla.service';
import { PlatformLocation } from '@angular/common';

@IonicPage({
  segment: 'cla-organization-app-modal'
})
@Component({
  selector: 'cla-organization-app-modal',
  templateUrl: 'cla-organization-app-modal.html'
})
export class ClaOrganizationAppModal {
  alert: any;
  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
    public alertCtrl: AlertController,
    public events: Events,
    private location: PlatformLocation
  ) {
    this.getDefaults();

    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
      if (this.alert) {
        this.alert.dismiss();
      }
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() { }

  getDefaults() { }

  // TODO: Do we want a call to cla that polls for the installation status?
  // UH YEA?


  openMessageDialog() {
    this.alert = this.alertCtrl.create({
      title: 'You are exiting EasyCLA and going to GitHub site. Please make sure you are already logged into GitHub so that you can install EasyCLA app',
      buttons: [
        {
          text: 'Ok',
          handler: () => {
            this.viewCtrl.dismiss();
            window.open(EnvConfig['gh-app-public-link'], '_blank');
          }
        },
      ]
    });
    this.alert.present();
  }

  openAppPage() {
    this.openMessageDialog()
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
