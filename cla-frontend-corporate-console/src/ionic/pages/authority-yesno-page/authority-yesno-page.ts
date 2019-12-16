// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'project/:projectId/yesno'
})
@Component({
  selector: 'authority-yesno-page',
  templateUrl: 'authority-yesno-page.html'
})
export class AuthorityYesnoPage {
  projectId: string;

  company: any;
  project: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.company = navParams.get('company');
  }

  getDefaults() {
    this.project = {
      project_name: ''
    };
  }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe((response) => {
      this.project = response;
    });
  }

  openClaCorporatePage() {
    this.navCtrl.push('ClaCorporatePage', {
      projectId: this.projectId,
      company: this.company
    });
  }

  openCollectAuthorityEmailModal() {
    let modal = this.modalCtrl.create('CollectAuthorityEmailModal', {
      projectId: this.projectId,
      companyId: this.company.company_id,
      projectName: this.project.project_name,
      companyName: this.company.company_name
    });
    modal.present();
  }
}
