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
  segment: 'cla/gerrit/project/:projectId/corporate'
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
    this.projectId = navParams.get('projectId');
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
    this.claService.getProjectCompanies(this.projectId).subscribe(response => {
      if (response) {
        this.companies = response;
      }
      this.loading.companies = false;
    });
  }

  getUserInfo() {
    // retrieve userInfo from auth0 service
    this.authService.getUserInfo().then(res => {
      this.claService.postOrGetUserForGerrit().subscribe(user => {
          this.userId = user.user_id;
          // get signatureIntent object, similar to the Github flow. 
          //this.postSignatureRequest();
      })
    })
  }


  
  postSignatureRequest() {
    let signatureRequest = {
      'project_id': this.projectId,
      'user_id': this.userId,
      'return_url_type': "Gerrit",
    };
    this.claService.postEmployeeSignatureRequest(signatureRequest).subscribe(response => {
      this.signature = response;
    });
  }

  openClaEmployeeCompanyConfirmPage(company) {
    let signatureRequest = {
      project_id: this.projectId,
      company_id: company.company_id
    };

    this.claService.postEmployeeSignatureRequest(signatureRequest).subscribe(response => {
      let errors = response.hasOwnProperty('errors');
      if (errors) {
        if (response.errors.hasOwnProperty('company_whitelist')) {
          // When the user is not whitelisted with the company: return {'errors': {'company_whitelist': 'User email (<email>) is not whitelisted for this company'}}
          this.openClaEmployeeCompanyTroubleshootPage(company);
          return;
        }
        if (response.errors.hasOwnProperty('missing_ccla')) {
          // When the company does NOT have a CCLA with the project: {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
          // The user shouldn't get here if they are using the console properly
          return;
        }
      } else {
        // No Errors, expect normal signature response
        this.signature = response;

      this.navCtrl.push('ClaEmployeeCompanyConfirmPage', {
        projectId: this.projectId
      });
      }
    });
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
    });
    modal.present();
    this.dismiss();
  }


  openClaEmployeeCompanyTroubleshootPage(company) {
    this.navCtrl.push('ClaEmployeeCompanyTroubleshootPage', {
      projectId: this.projectId,
      companyId: company.company_id,
    });
  }
  
}
