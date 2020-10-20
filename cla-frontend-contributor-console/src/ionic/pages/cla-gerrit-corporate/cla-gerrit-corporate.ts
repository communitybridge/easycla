// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { AuthService } from '../../services/auth.service';
import { Restricted } from '../../decorators/restricted';
import { generalConstants } from '../../constants/general';
import { EnvConfig } from '../../services/cla.env.utils';

@Restricted({
  roles: ['isAuthenticated']
})
@IonicPage({
  segment: 'cla/gerrit/project/:gerritId/corporate'
})
@Component({
  selector: 'cla-gerrit-corporate',
  templateUrl: 'cla-gerrit-corporate.html',
  providers: []
})
export class ClaGerritCorporatePage {
  loading: any;
  projectId: string;
  gerritId: string;
  userId: string;
  signature: string;
  companies: any;
  expanded: boolean = true;
  hasEnabledLFXHeader = EnvConfig['lfx-header-enabled'] === "true" ? true : false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    private claService: ClaService,
    private authService: AuthService,
  ) {
    this.gerritId = navParams.get('gerritId');
    this.getDefaults();
    localStorage.setItem('gerritId', this.gerritId);
    localStorage.setItem('gerritClaType', 'CCLA');
  }

  getDefaults() {
    this.loading = {
      companies: true
    };
    this.companies = [];
  }

  ngOnInit() {
    console.log(this.authService.loggedIn);
    if (!this.authService.loggedIn) {
      this.redirectToLogin();
    } else {
      this.getCompanies();
      this.getUserInfo();
      this.getProject();
    }
  }

  redirectToLogin() {
    this.navCtrl.setRoot('LoginPage');
  }

  getCompanies() {
    this.claService.getGerrit(this.gerritId).subscribe((gerrit) => {
      this.claService.getAllCompanies().subscribe((response) => {
        if (response) {
          this.companies = response;
        }
        this.loading.companies = false;
      });
    });
  }

  getUserInfo() {
    // retrieve userInfo from auth0 service
    this.claService.postOrGetUserForGerrit().subscribe((user) => {
      localStorage.setItem(generalConstants.USER_MODEL, JSON.stringify(user));
      this.userId = user.user_id;
    }, (error) => {
      // Got an auth error, redirect to the login
      console.log('Got an auth error, redirect to the login...');
      setTimeout(() => this.redirectToLogin());
    });
  }

  openClaEmployeeCompanyConfirmPage(company) {
    let data = {
      project_id: this.projectId,
      company_id: company.company_id,
      user_id: this.userId
    };
    this.claService.postCheckAndPreparedEmployeeSignature(data).subscribe((response) => {
      let errors = response.hasOwnProperty('errors');
      if (errors) {
        if (response.errors.hasOwnProperty('missing_ccla')) {
          // When the company does NOT have a CCLA with the project: {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
          this.openClaSendClaManagerEmailModal(company);
        }

        if (response.errors.hasOwnProperty('ccla_approval_list')) {
          // When the user is not whitelisted with the company: return {'errors': {'ccla_approval_list': 'No user email authorized for this ccla'}}
          this.openClaEmployeeCompanyTroubleshootPage(company);
          return;
        }
      } else {
        this.signature = response;

        this.navCtrl.push('ClaEmployeeCompanyConfirmPage', {
          projectId: this.projectId,
          signingType: 'Gerrit',
          userId: this.userId,
          companyId: company.company_id
        });
      }
    });
  }

  openClaSendClaManagerEmailModal(company) {
    let modal = this.modalCtrl.create('ClaSendClaManagerEmailModal', {
      projectId: this.projectId,
      userId: this.userId,
      companyId: company.company_id,
      authenticated: true
    });
    modal.present();
  }

  openClaNewCompanyModal() {
    let modal = this.modalCtrl.create('ClaNewCompanyModal', {
      projectId: this.projectId
    });
    modal.present();
  }

  openClaCompanyAdminYesnoModal() {
    let modal = this.modalCtrl.create('ClaCompanyAdminYesnoModal', {
      projectId: this.projectId,
      userId: this.userId,
      authenticated: true
    });
    modal.present();
  }

  openClaEmployeeCompanyTroubleshootPage(company) {
    this.navCtrl.push('ClaEmployeeCompanyTroubleshootPage', {
      projectId: this.projectId,
      repositoryId: '',
      userId: this.userId,
      companyId: company.company_id,
      gitService: 'Gerrit'
    });
  }

  getProject() {
    this.claService.getGerrit(this.gerritId).subscribe((response) => {
      if (response.errors) {
        console.error(response.errors);
        // Redirect to error page.
        this.redirectToLogin();
      } else {
        this.projectId = response.project_id;
        this.claService.getProjectWithAuthToken(response.project_id).subscribe((project) => {
          localStorage.setItem(generalConstants.PROJECT_MODEL, JSON.stringify(project));
        });
      }
    });
  }

  onClickToggle(hasExpanded) {
    this.expanded = hasExpanded;
  }
}
