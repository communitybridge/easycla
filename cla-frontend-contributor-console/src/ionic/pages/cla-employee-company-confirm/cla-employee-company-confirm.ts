// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CheckboxValidator } from '../../validators/checkbox';
import { ClaService } from '../../services/cla.service';
import { generalConstants } from '../../constants/general';

@IonicPage({
  segment: 'project/:projectId/user/:userId/employee/company/:companyId/confirm'
})
@Component({
  selector: 'cla-employee-company-confirm',
  templateUrl: 'cla-employee-company-confirm.html'
})
export class ClaEmployeeCompanyConfirmPage {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  signingType: string; // used to differentiate Github/Gerrit Users
  user: any;
  project: any;
  company: any;
  signature: any;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  errorMessage: string = null;
  helpDeskLink: URL = new URL(generalConstants.getHelpURL);

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.signingType = navParams.get('signingType');

    this.getDefaults();

    this.form = formBuilder.group({
      agree: [false, Validators.compose([CheckboxValidator.isChecked])]
    });
  }

  getDefaults() {
    this.project = {
      project_name: ''
    };
    this.company = {
      company_name: ''
    };
    this.user = {
      user_name: ''
    };
    this.errorMessage = null;
    this.currentlySubmitting = false;
  }

  ngOnInit() {
    this.getUser(this.userId);
    this.getProject(this.projectId);
    this.getCompany(this.companyId);
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe((response) => {
      this.user = response;
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe((response) => {
      this.project = response;
    });
  }

  getCompany(companyId) {
    this.claService.getCompany(companyId).subscribe((response) => {
      this.company = response;
    });
  }

  submit() {
    // Reset our status and error messages
    this.submitAttempt = true;
    this.errorMessage = null;
    this.currentlySubmitting = true;

    if (!this.form.valid) {
      this.currentlySubmitting = false;
      return;
    }

    let signatureRequest = {
      project_id: this.projectId,
      company_id: this.companyId,
      user_id: this.userId,
      return_url_type: this.signingType //"Gerrit" / "Github"
    };
    this.claService.postEmployeeSignatureRequest(signatureRequest).subscribe((response) => {
      this.currentlySubmitting = false;

      let errors = response.hasOwnProperty('errors');
      if (errors) {
        this.errorMessage = response.errors;

        if (response.errors.hasOwnProperty('ccla_approval_list')) {
          // When the user is not whitelisted with the company: return {'errors': {'ccla_approval_list': 'No user email authorized for this ccla'}}
          this.openClaEmployeeCompanyTroubleshootPage();
          return;
        }

        if (response.errors.hasOwnProperty('missing_ccla')) {
          // When the company does NOT have a CCLA with the project: {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
          // The user shouldn't get here if they are using the console properly
          return;
        }
      } else {
        // No Errors, expect normal signature response
        this.errorMessage = null;
        this.signature = response;
        this.openClaNextStepModal();
      }
    });
  }

  openClaNextStepModal() {
    let modal = this.modalCtrl.create('ClaNextStepModal', {
      projectId: this.projectId,
      userId: this.userId,
      project: this.project,
      signature: this.signature,
      signingType: this.signingType
    });
    modal.present();
  }

  openClaEmployeeCompanyTroubleshootPage() {
    this.navCtrl.push('ClaEmployeeCompanyTroubleshootPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: this.companyId
    });
  }
}
