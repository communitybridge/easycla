import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController, NavParams, ViewController, } from 'ionic-angular';
import { ClaService } from 'cla-service';

@IonicPage({
  segment: 'cla-contract-companies-modal'
})
@Component({
  selector: 'cla-contract-companies-modal',
  templateUrl: 'cla-contract-companies-modal.html',
})
export class ClaContractCompaniesModal {

  claProjectId: string;
  companies: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private claService: ClaService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {
    this.getProjectCompanies();
  }

  getDefaults() {
    this.claProjectId = this.navParams.get('claProjectId');
    this.companies = [];
  }

  getProjectCompanies() {
    this.claService.getProjectCompanies(this.claProjectId).subscribe(companies => {
      this.companies = companies;
    });
  }

  openClaCorporateMemberOptionsModal() {
    let modal = this.modalCtrl.create('ClaCorporateMemberOptionsModal');
    modal.present();
  }

  openClaCorporateManagementPage(companyId) {
    console.log(companyId);
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }



}
