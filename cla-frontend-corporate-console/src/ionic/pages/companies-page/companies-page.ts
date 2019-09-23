// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from "@angular/core";
import { NavController, ModalController, IonicPage } from "ionic-angular";
import { ClaService } from "../../services/cla.service";
import { ClaCompanyModel } from "../../models/cla-company";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";
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
    this.userId = localStorage.getItem("userid")
    this.companies = [];
    this.columns = [
      {prop: 'CompanyName'},
      {prop: 'Status'},
    ];
    this.rows = [];
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
    this.claService.getCompanies().subscribe(response => {
      let company = response[0];
      company && this.getCompaniesByUserManagerWithInvites(company.company_manager_id);
    });
  }

  getCompaniesByUserManagerWithInvites(userId) {
    this.claService.getCompaniesByUserManagerWithInvites(userId).subscribe((companies) => {
      this.loading.companies = false;
      this.rows = this.mapCompanies(companies['companies-with-invites']);
      console.log(this.rows, 'rows')
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
    console.log(event)
    // this.openProjectPage(event.selected[0].ProjectID);
  }

  mapCompanies(companies) {
    let rows = [];
    for (let company of companies) {
      rows.push({
        companyID: company.companyID,
        CompanyName: company.companyName,
        Status: company.status,
      });
    }
    return rows;
  }
}
