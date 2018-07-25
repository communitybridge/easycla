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
    this.navCtrl.push('ClaEmployeeCompanyConfirmPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: company.company_id,
    });
  }

  openClaNewCompanyModal() {
    let modal = this.modalCtrl.create('ClaNewCompanyModal', {
      projectId: this.projectId,
    });
    modal.present();
  }

}
