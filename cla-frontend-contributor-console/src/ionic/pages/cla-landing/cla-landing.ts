// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, IonicPage, ModalController, NavParams, } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId'
})
@Component({
  selector: 'cla-landing',
  templateUrl: 'cla-landing.html'
})
export class ClaLandingPage {
  loading: any;
  projectId: string;
  userId: string;

  user: any;
  project: any;

  hasIndividualCla: boolean;
  hasCorporateCla: boolean;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      individualDoc: true,
      corporateDoc: true,
    }
    this.project = {
      project_name: "",
    }

    this.hasCorporateCla = false;
    this.hasIndividualCla = false;
  }

  ngOnInit() {
    this.getUser(this.userId);
    this.getProject(this.projectId);
    this.getProjectDocuments();
  }

  openClaIndividualPage() {
    // send to the individual cla page which will give directions and redirect
    this.navCtrl.push('ClaIndividualPage', {
      projectId: this.projectId,
      userId: this.userId,
    });
  }

  openClaIndividualEmployeeModal() {
    let modal = this.modalCtrl.create('ClaSelectCompanyModal', {
      projectId: this.projectId,
      userId: this.userId,
    });
    modal.present();
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.user = response;
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
    });
  }

  getProjectDocuments() {
    this.claService.getProjectDocument(this.projectId, 'individual').subscribe(response => {
      if (!response.hasOwnProperty('errors')) {
        this.hasIndividualCla = true;
      }
      this.loading.individualDoc = false;
    });
    this.claService.getProjectDocument(this.projectId, 'corporate').subscribe(response => {
      if (!response.hasOwnProperty('errors')) {
        this.hasCorporateCla = true;
      }
      this.loading.corporateDoc = false;
    });
  }

}
