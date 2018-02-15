import { Component } from '@angular/core';
import { NavController, ModalController, IonicPage } from 'ionic-angular';
import { ClaService } from 'cla-service';


@IonicPage({
  segment: 'companies'
})
@Component({
  selector: 'companies-page',
  templateUrl: 'companies-page.html'
})
export class CompaniesPage {
  loading: any;
  companies: any;

  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    console.log('companies page');
    this.loading = {
      companies: true,
    };

    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
  }

  openCompanyModel() {
    let modal = this.modalCtrl.create('AddCompanyModal', {});
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompanies();
    });
    modal.present();
  }

  getCompanies() {
    console.log('get companies called');
    this.claService.getCompanies().subscribe(response => {
        this.companies = response;
        console.log(this.companies);
        this.loading.companies = false;
    });
  }

  viewCompany(companyId){
    this.navCtrl.setRoot('CompanyPage', {
      companyId: companyId
    });
  }

}
