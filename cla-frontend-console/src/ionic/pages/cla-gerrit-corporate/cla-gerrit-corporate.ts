// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { RolesService } from "../../services/roles.service";
import { AuthService } from "../../services/auth.service";
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { Restricted } from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: 'cla/gerrit/project/:gerritId/corporate'
})
@Component({
  selector: 'cla-gerrit-corporate',
  templateUrl: 'cla-gerrit-corporate.html',
  providers: [
  ]
})
export class ClaGerritCorporatePage {
  loading: any;
  projectId: string;
  gerritId: string;
  userId: string;

  signature: string;
  
  companies: any;
  

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    private claService: ClaService,
    private rolesService: RolesService,
    private authService: AuthService,
    private keycloak: KeycloakService,
  ) {
    this.gerritId = navParams.get('gerritId');
    this.getDefaults();
    localStorage.setItem("gerritId", this.gerritId);
    localStorage.setItem("gerritClaType", "CCLA");
  }

  getDefaults() {
    this.loading = {
      companies: true,
    };
    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
    this.getUserInfo();
  }

  ionViewCanEnter(){
    if(!this.authService.isAuthenticated){
      setTimeout(()=>this.navCtrl.setRoot('LoginPage'))
    }
    return this.authService.isAuthenticated
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  getCompanies() {
    this.claService.getGerrit(this.gerritId).subscribe(gerrit => {
      this.projectId = gerrit.project_id;
      this.claService.getProjectCompanies(gerrit.project_id).subscribe(response => {
        if (response) {
          this.companies = response;
        }
        this.loading.companies = false;
      });
    });
  }

  getUserInfo() {
    // retrieve userInfo from auth0 service
      this.claService.postOrGetUserForGerrit().subscribe(user => {
          this.userId = user.user_id;
          console.log(this.userId);
          // get signatureIntent object, similar to the Github flow. 
          //this.postSignatureRequest();
      })
  }

  openClaEmployeeCompanyConfirmPage(company) {
    let data = {
      project_id: this.projectId,
      company_id: company.company_id,
      user_id: this.userId
    };

    this.claService.postCheckedAndPreparedEmployeeSignature(data).subscribe(response => {
      let errors = response.hasOwnProperty('errors');
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
        this.signature = response;

      this.navCtrl.push('ClaEmployeeCompanyConfirmPage', {
        projectId: this.projectId,
        signingType: "Gerrit",
        userId: this.userId,
        companyId: company.company_id,
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
      projectId: this.projectId,
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
      repositoryId: "",
      userId: this.userId,
      companyId: company.company_id,
      gitService: 'Gerrit'
    });
  }
}
