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

  individualClaAgreement: string;
  corporateClaAgreement: string;

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
    this.individualClaAgreement = "";
    this.corporateClaAgreement = ""
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
      if (response.document_s3_url) {
        this.individualClaAgreement = response.document_s3_url;
      }
      this.loading.individualDoc = false;
    });
    this.claService.getProjectDocument(this.projectId, 'corporate').subscribe(response => {
      if (!response.hasOwnProperty('errors')) {
        this.hasCorporateCla = true;
      }
      if(response.document_s3_url) {
        this.corporateClaAgreement = response.document_s3_url;
      }
      this.loading.corporateDoc = false;
    });
  }

  getUserFriendlyID() {
    if (this.user == null) {
      return "user is undefined";
    }

    if (this.user.user_github_username != null) {
      return this.user.user_github_username + " Github user id";
    } else if (this.user.lf_username != null) {
        return this.user.lf_username + " LF user id";
    } else if (this.user.lf_email != null) {
      return this.user.lf_email + " LF email";
    } else {
      return "no user information";
    }
  }

}
