// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/individual'
})
@Component({
  selector: 'cla-individual',
  templateUrl: 'cla-individual.html'
})
export class ClaIndividualPage {
  projectId: string;
  userId: string;

  user: any;
  project: any;
  signatureIntent: any;
  activeSignatures: boolean = true; // we assume true until otherwise
  signature: any;
  loadingSignature: boolean = true;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
  }

  getDefaults() {
    this.project = {
      project_name: "",
    };
    this.signature = {
      sign_url: "",
    };
  }

  ngOnInit() {
    this.getUser(this.userId);
    this.getProject(this.projectId);
    this.getUserSignatureIntent(this.userId);
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

  getUserSignatureIntent(userId)  {
    this.loadingSignature = true;
    this.claService.getUserSignatureIntent(userId).subscribe(response => {
      this.signatureIntent = response;
      if(this.signatureIntent !== null) {
        this.postSignatureRequest();
      } else {
        this.activeSignatures = false;
      }
      this.loadingSignature = false;
    });
  }

  postSignatureRequest() {
    let signatureRequest = {
      'project_id': this.projectId,
      'user_id': this.userId,
      // TODO: Switch this to intermediary loading screen as docusign postback has delay
      'return_url_type': "Github",
      'return_url': this.signatureIntent.return_url,
    };

    this.claService.postIndividualSignatureRequest(signatureRequest).subscribe(response => {
      // returns {
      //   user_id:
      //   signature_id:
      //   project_id:
      //   sign_url: docusign.com/some-docusign-url
      // }
      console.log(response);
      this.signature = response;
    });
  }

  openClaAgreement() {
    if (!this.signature.sign_url) {
      // Can't open agreement if we don't have a sign_url yet
      return;
    }
    window.open(this.signature.sign_url, '_blank');
  }

}
