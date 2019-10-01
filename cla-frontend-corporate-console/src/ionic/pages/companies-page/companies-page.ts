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
  userEmail: string;
  userName: string;

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
    this.userEmail = localStorage.getItem("user_email");
    this.userName = localStorage.getItem("user_name");
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

  /**
   * Get the list of companies for this user - if any
   */
  getCompanies() {
    this.loading.companies = true;
    this.claService.getUserByUserName(this.userId).subscribe(
      response => {
        console.log(response);
        // We need the user's unique ID - grab it from the first record
        this.getCompaniesByUserManagerWithInvites(response.userID);
      },
      exception => {
        // Typically get this if the user is not found in our database
        this.loading.companies = false;
        this.rows = [];

        console.log('Received exception:');
        console.log(exception);
        if (exception.status === 404) {
          // Create the user if it does't exist
          console.log('Creating user record...');
          const user = {
            'lfUsername': this.userId,
            'username': this.userName,
            'lfEmail': this.userEmail,
          };
          this.claService.createUserV3(user).subscribe(
            response => {
              console.log('Success creating user record: ');
              console.log(response);
            },
            exception => {
              console.log('Error creating user record: ');
              console.log(exception);
            });
        }

      });
  }

  /**
   * Fetch the list of companies and company managers - update the companies table view
   *
   * @param userId the username/id of the logged in user
   */
  getCompaniesByUserManagerWithInvites(userId) {
    this.claService.getCompaniesByUserManagerWithInvites(userId).subscribe((companies) => {
        this.loading.companies = false;
        if (companies['companies-with-invites']) {
          this.rows = this.mapCompanies(companies['companies-with-invites']);
        } else{
          this.rows = [];
        }
      },
      exception => {
        this.loading.companies = false;
        console.log("Exception while calling: getCompaniesByUserManagerWithInvites() for userId: " + userId);
        console.log(exception);
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
