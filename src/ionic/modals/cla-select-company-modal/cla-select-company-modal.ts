import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company'
})
@Component({
  selector: 'cla-select-company-modal',
  templateUrl: 'cla-select-company-modal.html',
  providers: [
    CincoService,
  ]
})
export class ClaSelectCompanyModal {
  projectId: string;
  repositoryId: string;
  userId: string;

  companies: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    private modalCtrl: ModalController,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
  }

  getDefaults() {
    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  getCompanies() {
    this.companies = [
      {
        name: "Alphabet (Google)",
        id: "0000000001",
      },
      {
        name: "Apple",
        id: "0000000002",
      },
      {
        name: "Company Name",
        id: "0000000003",
      },
      {
        name: "Company Name",
        id: "0000000004",
      },
      {
        name: "Company Name",
        id: "0000000005",
      },
      {
        name: "Company Name",
        id: "0000000006",
      },
      {
        name: "Company Name",
        id: "0000000007",
      },
      {
        name: "Company Name",
        id: "0000000008",
      },
      {
        name: "Company Name",
        id: "0000000009",
      },
      {
        name: "Company Name",
        id: "0000000010",
      },
      {
        name: "Company Name",
        id: "0000000011",
      },
    ];
  }

  openClaEmployeeCompanyPage(company) {
    this.navCtrl.push('ClaEmployeeCompanyPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: company.id,
    });
  }

  openClaNewCompanyModal() {
    let modal = this.modalCtrl.create('ClaNewCompanyModal', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
    });
    modal.present();
  }

}
