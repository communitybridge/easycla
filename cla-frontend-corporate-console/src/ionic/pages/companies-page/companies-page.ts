// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {IonicPage, ModalController, NavController} from "ionic-angular";
import {ClaService} from "../../services/cla.service";
import {RolesService} from "../../services/roles.service";
import {Restricted} from "../../decorators/restricted";
import {ColumnMode, SelectionType, SortType} from "@swimlane/ngx-datatable";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: "companies"
})
@Component({
  selector: "companies-page",
  templateUrl: "companies-page.html"
})
export class CompaniesPage {
  loading: any;
  companies: any;
  userId: string;
  manager: string;
  columns: any[];
  rows: any[];

  ColumnMode = ColumnMode;
  SelectionType = SelectionType;
  SortType = SortType;

  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private rolesService: RolesService // for @Restricted
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companies: true
    };
    this.userId = localStorage.getItem("userid");
    this.companies = [];
    this.columns = [
      {prop: 'CompanyName'},
      {prop: 'Status'}
    ];
  }

  ngOnInit() {
    this.getCompanies();
  }

  openCompanyModal() {
    let modal = this.modalCtrl.create("AddCompanyModal", {});
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompanies();
    });
    modal.present();
  }

  getCompanies() {
    this.loading.companies = true;
    this.claService.getUserByUserName(this.userId).subscribe(response => {
      if (response != null) {
        // We need the user's unique ID - grab it from the first record
        this.getCompaniesByUserManagerWithInvites(response.userID);
      } else {
        this.loading.companies = false;
      }
    });
  }

  getCompaniesByUserManagerWithInvites(userId) {
    this.claService.getCompaniesByUserManagerWithInvites(userId).subscribe((companies) => {
      this.loading.companies = false;
      this.rows = this.mapCompanies(companies['companies-with-invites']);
    })
  }

  getPendingUserInvite(companyId, userId) {
    this.claService.getPendingUserInvite(companyId, userId).subscribe((response) => {
    })
  }

  viewCompany(companyId) {
    this.navCtrl.setRoot("CompanyPage", {
      companyId: companyId
    });
  }

  onSelect(event) {
    let company = event.selected[0]
    if (company.Status === "Joined") {
      this.viewCompany(company.companyID);
    }
  }

  mapCompanies(companies) {
    let rows = [];
    for (let company of companies) {
      rows.push({
        companyID: company.companyID,
        CompanyName: company.companyName,
        Status: company.status
      });
    }
    return rows;
  }
}
