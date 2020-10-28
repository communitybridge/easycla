// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController, NavParams } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { generalConstants } from '../../constants/general';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId'
})
@Component({
  selector: 'cla-landing',
  templateUrl: 'cla-landing.html'
})
export class ClaLandingPage {
  projectId: string;
  userId: string;
  user: any;
  project: any;
  expanded: boolean = true;
  hasEnabledLFXHeader = EnvConfig['lfx-header-enabled'] === "true" ? true : false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.getDefaults();
  }

  getDefaults() {
    this.project = {
      project_name: ''
    };
  }

  ngOnInit() {
    this.getUser();
    this.getProject();
    localStorage.removeItem('gerritClaType');
  }

  openClaIndividualPage() {
    // send to the individual cla page which will give directions and redirect
    this.navCtrl.push('ClaIndividualPage', {
      projectId: this.projectId,
      userId: this.userId
    });
  }

  openClaIndividualEmployeeModal() {
    let modal = this.modalCtrl.create('ClaSelectCompanyModal', {
      projectId: this.projectId,
      userId: this.userId
    });
    modal.present();
  }

  getUser() {
    this.claService.getUser(this.userId).subscribe((response) => {
      localStorage.setItem(generalConstants.USER_MODEL, JSON.stringify(response));
      this.user = response;
    });
  }

  getProject() {
    this.claService.getProject(this.projectId).subscribe((response) => {
      localStorage.setItem(generalConstants.PROJECT_MODEL, JSON.stringify(response));
      this.project = response;
    });
  }

  onClickToggle(hasExpanded) {
    this.expanded = hasExpanded;
  }
}
