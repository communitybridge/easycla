import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
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

  signature: string;
  
  companies: any;
  

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
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

    let signatureRequest = {
      project_id: this.projectId,
      company_id: company.company_id,
      user_id: this.userId,
    };

    this.claService.postEmployeeSignatureRequest(signatureRequest).subscribe(response => {
      let errors = response.hasOwnProperty('errors');
      this.selectCompanyModalActive = false;
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
          projectId: this.projectId,
          repositoryId: this.repositoryId,
          userId: this.userId,
          companyId: company.company_id,
          signingType: "Github",
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
      userId: this.userId
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
    });
  }
  
}
