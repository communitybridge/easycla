import { Component, ViewChild } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { CincoService } from '../../services/cinco.service'

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
    private cincoService: CincoService,
    private keycloak: KeycloakService
  ) {
    this.getDefaults();
  }

  getDefaults(){
    this.loading = {
      companies: true,
    };

    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
  }

  ionViewCanEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.setRoot('LoginPage');
      this.navCtrl.popToRoot();
    }
    return this.keycloak.authenticated();
  }

  ionViewWillEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.push('LoginPage');
    }
  }

  getCompanies() {
    // this.cincoclaService.getCompaniesByUser().subscribe(response => {
    //     this.companies = response;
    //     this.loading.companies = false;
    // });
    this.companies = [
      {
        company_name: "abc",
        company_id: "123",
      }
    ];
    this.loading.companies = false;
  }

  viewCompany(companyId){
    this.navCtrl.setRoot('CompanyPage', {
      companyId: companyId
    });
  }

}
