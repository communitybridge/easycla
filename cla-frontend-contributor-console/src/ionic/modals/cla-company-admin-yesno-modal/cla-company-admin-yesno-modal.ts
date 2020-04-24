// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/admin-yesno'
})
@Component({
  selector: 'cla-company-admin-yesno-modal',
  templateUrl: 'cla-company-admin-yesno-modal.html'
})
export class ClaCompanyAdminYesnoModal {
  projectId: string;
  companyId: string;
  companyName: string;
  userId: string;
  authenticated: boolean; //true if coming from gerrit/corporate
  consoleLink: string;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController
  ) {
    this.projectId = navParams.get('projectId');
    // May be empty
    this.companyId = navParams.get('companyId') || '';
    // May be empty
    this.companyName = navParams.get('companyName') || '';
    this.userId = navParams.get('userId');
    this.authenticated = navParams.get('authenticated');
    this.getDefaults();
  }

  getDefaults() {
    this.consoleLink = EnvConfig['corp-console-link'];
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  openCompanyAdminConsoleLink() {
    window.open(this.consoleLink, '_blank');
  }

  openCompanyAdminSendEmail() {
    let modal = this.modalCtrl.create('ClaCompanyAdminSendEmailModal', {
      projectId: this.projectId,
      companyId: this.companyId,
      companyName: this.companyName,
      userId: this.userId,
      authenticated: this.authenticated
    });
    modal.present();
    this.dismiss();
  }
}
