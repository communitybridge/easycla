import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { ClaService } from 'cla-service';
import { ClaCompanyModel } from '../../models/cla-company';
import { ClaUserModel } from '../../models/cla-user';
import { ClaSignatureModel } from '../../models/cla-signature';
// import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';
// import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'company/:companyId'
})
@Component({
  selector: 'company-page',
  templateUrl: 'company-page.html',
})
export class CompanyPage {
  companyId: string;
  company: ClaCompanyModel;
  manager: ClaUserModel;
  companySignatures: ClaSignatureModel[];
  projects: any;
  loading: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public modalCtrl: ModalController,
  ) {
    this.companyId = navParams.get('companyId');
    console.log(this.companyId);
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companySignatures: true,
    };
    this.company = new ClaCompanyModel();
    this.projects = {};
  }

  ngOnInit() {
    this.getCompany();
    this.getCompanySignatures();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      console.log(response);
      this.company = response;
      this.getUser(this.company.company_manager_id);
    });
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      console.log('getUser: ' + userId);
      console.log(response);
      this.manager = response;
    });
  }

  getCompanySignatures() {
    this.claService.getCompanySignatures(this.companyId).subscribe(response => {
      console.log('company signatures');
      console.log(response);
      this.companySignatures = response;
      for (let signature of this.companySignatures) {
        this.getProject(signature.signature_project_id);
      }
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      console.log('project: ' + projectId);
      console.log(response);
      this.projects[projectId] = response;
      console.log(this.projects);
    });
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      companyId: this.companyId,
      projectId: projectId,
    });
  }
  
  openCompanyModal() {
    console.log('company page');
    console.log(this.company);
    let modal = this.modalCtrl.create('AddCompanyModal', {
      company: this.company,
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

}
