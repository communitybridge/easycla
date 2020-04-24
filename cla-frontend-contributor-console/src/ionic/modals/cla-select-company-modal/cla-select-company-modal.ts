// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company'
})
@Component({
  selector: 'cla-select-company-modal',
  templateUrl: 'cla-select-company-modal.html',
  providers: []
})
export class ClaSelectCompanyModal {
  form: FormGroup;
  loading: any;
  projectId: string;
  repositoryId: string;
  userId: string;
  selectCompanyModalActive: boolean = false;
  authenticated: boolean;
  signature: string;
  companies: any;
  companiesFiltered: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    public formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.authenticated = navParams.get('authenticated');
    this.getDefaults();
    this.form = formBuilder.group({
      search: ['', Validators.compose([Validators.required])]
    });
  }

  getDefaults() {
    this.loading = {
      companies: true,
      activateSpinner: false
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
    this.loading.companies = true;
    this.claService.getAllCompanies().subscribe((response) => {
      this.loading.companies = false;
      if (response) {
        // Cleanup - Remove any companies that don't have a name
        this.companies = response.filter((company) => {
          return company.company_name && company.company_name.trim().length > 0;
        });

        // Reset our filtered search
        this.form.value.search = '';
        this.companiesFiltered = this.companies;
      }
    });
  }

  openClaEmployeeCompanyConfirmPage(company) {
    // set loading spinner to true when a company is selected
    this.loading.activateSpinner = true;
    if (this.selectCompanyModalActive) {
      return false;
    }
    this.selectCompanyModalActive = true;

    let data = {
      project_id: this.projectId,
      company_id: company.company_id,
      user_id: this.userId
    };

    this.claService.postCheckedAndPreparedEmployeeSignature(data).subscribe((response) => {
      /*
      Before an employee begins the signing process, ensure that
      1. The given project, company, and user exists
      2. The company signatory has signed the CCLA for their company.
      3. The user is included as part of the approved list of the CCLA that the company signed.
      the CLA service will throw an error if any of the above is false.
      */
      this.loading.activateSpinner = false;
      let errors = response.hasOwnProperty('errors');
      console.log(`errors: ${errors}`);
      this.selectCompanyModalActive = false;
      if (errors) {
        if (response.errors.hasOwnProperty('missing_ccla')) {
          // When the company does NOT have a CCLA with the project: {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
          console.log(`errors - missing_ccla: ${response}`);
          //this.openClaSendClaManagerEmailModal(company);
          this.openClaCompanyAdminYesnoModal(company);
        }

        if (response.errors.hasOwnProperty('ccla_whitelist')) {
          console.log(`errors - ccla_whitelist: ${response}`);
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
          signingType: 'Github'
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

  openClaCompanyAdminYesnoModal(company) {
    let modal = this.modalCtrl.create('ClaCompanyAdminYesnoModal', {
      projectId: this.projectId,
      companyId: company == null ? '' : company.company_id,
      companyName: company == null ? '' : company.company_name,
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

  onSearch() {
    const searchTerm = this.form.value.search;
    if (searchTerm === '') {
      this.companiesFiltered = this.companies;
    } else {
      this.companiesFiltered = this.companies.filter((a) => {
        return a.company_name.toLowerCase().includes(searchTerm.toLowerCase());
      });
    }
  }
}
