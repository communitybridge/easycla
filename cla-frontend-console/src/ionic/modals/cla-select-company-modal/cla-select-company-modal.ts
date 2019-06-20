// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage, AlertController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company'
})
@Component({
  selector: 'cla-select-company-modal',
  templateUrl: 'cla-select-company-modal.html',
  providers: [
  ]
})
export class ClaSelectCompanyModal {
  loading: any;
  projectId: string;
  repositoryId: string;
  userId: string;
  selectCompanyModalActive: boolean = false;
  authenticated: boolean;

  signature: string;

  companies: any;


  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    public alertCtrl: AlertController,
    private claService: ClaService,
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.authenticated = navParams.get('authenticated');
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companies: true,
    };
    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  getCompanies() {
    this.claService.getProjectCompanies(this.projectId).subscribe(response => {
      if (response) {
        this.companies = response;
      }
      this.loading.companies = false;
    });
  }

  openClaEmployeeCompanyConfirmPage(company) {
    if(this.selectCompanyModalActive){
      return false;
    }
    this.selectCompanyModalActive = true;

    let data = {
      project_id: this.projectId,
      company_id: company.company_id,
      user_id: this.userId
    };

    this.claService.postCheckedAndPreparedEmployeeSignature(data).subscribe(response => {
      /* 
      Before an employee begins the signing process, ensure that
      1. The given project, company, and user exists 
      2. The company signatory has signed the CCLA for their company. 
      3. The user is included as part of the whitelist of the CCLA that the company signed. 
      the CLA service will throw an error if any of the above is false. 
      */
      let errors = response.hasOwnProperty('errors');
      this.selectCompanyModalActive = false;
      if (errors) {

        if (response.errors.hasOwnProperty('missing_ccla')) {
          // When the company does NOT have a CCLA with the project: {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
          this.openClaSendClaManagerEmailModal(company);
        }

        if (response.errors.hasOwnProperty('ccla_whitelist')) {
          // When the user is not whitelisted with the company: return {'errors': {'ccla_whitelist': 'No user email whitelisted for this ccla'}}
          this.openClaEmployeeCompanyTroubleshootPage(company);
          return;
        }

      } else {
        // No Errors, expect normal signature response
        this.signature = response;

        this.navCtrl.push('ClaEmployeeCompanyConfirmPage', {
          projectId: this.projectId,
          repositoryId: this.repositoryId,
          userId: this.userId,
          companyId: company.company_id,
          signingType: "Github"
        });
      }
    });
  }

  openClaSendClaManagerEmailModal(company) {
    let modal = this.modalCtrl.create('ClaSendClaManagerEmailModal', {
      projectId: this.projectId,
      userId: this.userId,
      companyId: company.company_id,
      authenticated: this.authenticated
    });
    modal.present();
  }

  openClaCompanyAdminYesnoModal() {
    let modal = this.modalCtrl.create('ClaCompanyAdminYesnoModal', {
      projectId: this.projectId,
      userId: this.userId,
      authenticated: false // Github users are not authenticated.
    });
    modal.present();
    this.dismiss();
  }


  openClaEmployeeCompanyTroubleshootPage(company) {
    this.navCtrl.push('ClaEmployeeCompanyTroubleshootPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: company.company_id,
      gitService: 'GitHub'
    });
  }

}
