// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company/:companyId/contact/sent'
})
@Component({
  selector: 'cla-message-sent',
  templateUrl: 'cla-message-sent.html'
})
export class ClaMessageSentPage {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;

  project: any;
  company: any;

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
  ) {

    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
  }

  getDefaults(){

  }

  ngOnInit(){
    this.getProject();
    this.getCompany();
  }

  getProject() {
    this.project = {
      id: '0000000001',
      name: 'Project Name',
      logoRef: 'https://dummyimage.com/225x102/d8d8d8/242424.png&text=Project+Logo',
    };
  }

  getCompany() {
    this.company = {
      name: 'Company Name',
      id: '0000000001',
    };
  }

  openClaIndividualPage() {
    // send to the individual cla page which will give directions and redirect
    this.navCtrl.push('ClaIndividualPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
    });
  }

}
