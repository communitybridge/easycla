import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company/:companyId/troubleshoot'
})
@Component({
  selector: 'cla-employee-company-troubleshoot',
  templateUrl: 'cla-employee-company-troubleshoot.html'
})
export class ClaEmployeeCompanyTroubleshootPage {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  authenticated: boolean;

  project: any;
  company: any;
  gitService: string;

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.authenticated = navParams.get('authenticated'); 
  }

  getDefaults() {

  }

  ngOnInit() {
    this.getProject();
    this.getCompany();
    this.getGitService();
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

  getGitService() {
    // GitHub, GitLab, Gerrit
    this.gitService = 'GitHub';
  }

  openGitServiceEmailSettings() {
    window.open("https://github.com/settings/emails", "_blank");
  }

  openClaEmployeeRequestAccessModal() {
    let modal = this.modalCtrl.create('ClaEmployeeRequestAccessModal', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: this.companyId,
      authenticated: this.authenticated
    });
    modal.present();
  }

}
