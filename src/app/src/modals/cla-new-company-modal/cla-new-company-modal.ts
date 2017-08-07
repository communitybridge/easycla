import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company'
})
@Component({
  selector: 'cla-new-company-modal',
  templateUrl: 'cla-new-company-modal.html',
  providers: [
    CincoService,
  ]
})
export class ClaNewCompanyModal {
  projectId: string;
  repositoryId: string;
  userId: string;

  cclaLink: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
  }

  getDefaults() {
    this.cclaLink = "";
  }

  ngOnInit() {
    this.getCclaLink();
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  getCclaLink() {
    this.cclaLink = "https://docusign.com/some/really/long/link/that/is/long/enough/to/push/it/out/of/the/input/viewport";
  }

  openClaEmployeeCompanyPage(company) {
    this.navCtrl.push('ClaEmployeeCompanyPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: company.id,
    });
  }

  openCclaLink() {
    window.open(this.cclaLink, '_blank');
  }


  // sortContacts(prop) {
  //   this.sortService.toggleSort(
  //     this.sort,
  //     prop,
  //     this.organizationContacts,
  //   );
  // }

}
