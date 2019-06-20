// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company/:companyId/troubleshoot'
})
@Component({
  selector: 'cla-employee-company-troubleshoot',
  templateUrl: 'cla-employee-company-troubleshoot.html'
})
export class ClaEmployeeCompanyTroubleshootPage {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  authenticated: boolean;

  project: any;
  company: any;
  gitService: string;

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.gitService = navParams.get('gitService');
    this.authenticated = navParams.get('authenticated'); 
  }

  getDefaults() {
    this.project = {
      project_name: '',
      logoUrl: '',
    };
    this.company = {
      company_name: '',
    };
  }

  ngOnInit() {
    this.getProject(this.projectId);
    this.getCompany(this.companyId);
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
    });
  }

  getCompany(companyId) {
    this.claService.getCompany(companyId).subscribe(response => {
      this.company = response;
    });
  }

  openGitServiceEmailSettings() {
    window.open("https://github.com/settings/emails", "_blank");
  }

  openClaEmployeeRequestAccessModal() {
    let modal = this.modalCtrl.create('ClaEmployeeRequestAccessModal', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: this.companyId,
      authenticated: this.authenticated
    });
    modal.present();
  }

}
