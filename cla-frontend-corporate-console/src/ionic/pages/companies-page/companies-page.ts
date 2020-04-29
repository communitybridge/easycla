// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { Restricted } from '../../decorators/restricted';
import { RolesService } from '../../services/roles.service';

@Restricted({
  roles: ['isAuthenticated']
})
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
  userId: string;
  userEmail: string;
  userName: string;
  companyId: string;
  rows: any[];
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

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
    this.userId = localStorage.getItem('userid');
    this.userEmail = localStorage.getItem('user_email');
    this.userName = localStorage.getItem('user_name');
    this.companies = [];
  }

  ngOnInit() {
    this.getCompanies();
  }

  /**
   * Get the list of companies for this user - if any
   */
  getCompanies() {
    this.loading.companies = true;
    this.claService.getUserByUserName(this.userId).subscribe(
      (response) => {
        // We need the user's unique ID - grab it from the first record
        this.getCompaniesByUserManagerWithInvites(response.userID);
      },
      (exception) => {
        // Typically get this if the user is not found in our database
        this.loading.companies = false;
        this.rows = [];
        if (exception.status === 404) {
          // Create the user if it does't exist
          const user = {
            lfUsername: this.userId,
            username: this.userName,
            lfEmail: this.userEmail
          };
          this.claService.createUserV3(user).subscribe(
            (response) => {
              // Success creating user record
            },
            (exception) => {
              // Error creating user record
            }
          );
        }
      }
    );
  }

  /**
   * Fetch the list of companies and company managers - update the companies table view
   *
   * @param userId the username/id of the logged in user
   */
  getCompaniesByUserManagerWithInvites(userId) {
    this.claService.getCompaniesByUserManagerWithInvites(userId).subscribe(
      (companies) => {
        this.loading.companies = false;
        if (companies['companies-with-invites']) {
          this.rows = this.mapCompanies(companies['companies-with-invites']);
        } else {
          this.rows = [];
        }
      },
      (exception) => {
        this.loading.companies = false;
      }
    );
  }

  viewCompany(companyId, status) {
    if (status !== 'Pending Approval') {
      this.navCtrl.setRoot('CompanyPage', {
        companyId: companyId
      });
    }
  }

  mapCompanies(companies) {
    let rows = [];
    companies = this.sortData(companies);
    companies = this.removeDuplicates(companies);
    for (let company of companies) {
      rows.push({
        CompanyID: company.companyID,
        CompanyName: company.companyName,
        Status: company.status,
        ProjectName: ''
      });
    }
    return rows;
  }

  sortData(companies: any[]) {
    let joinedCompanies = companies.filter(company => company.status !== 'pending')
    let requestCompanies = companies.filter(company => company.status === 'pending')
    joinedCompanies = joinedCompanies.sort((a, b) => {
      return a.companyName.toLowerCase().localeCompare(b.companyName.toLowerCase());
    });
    requestCompanies = requestCompanies.sort((a, b) => {
      return a.companyName.toLowerCase().localeCompare(b.companyName.toLowerCase());
    });
    return joinedCompanies.concat(requestCompanies);
  }

  removeDuplicates(companies: any[]) {
    let seenCompanies = {};
    return companies.filter(company => {
      if (seenCompanies[company.companyID] == null) {
        seenCompanies[company.companyID] = true;
        return true; // unique, pass filter
      }
      return false; // duplicate, fail filter
    });
  }

  openSelectCompany() {
    if (!this.loading.companies) {
      //console.log(this.rows);
      let modal = this.modalCtrl.create('AddCompanyModal', {
        // For this modal, we share the list of companies where the user is either pending acceptance, approved/joined.
        // From this, the modal dialog will filter the list so that these companies are not shown. We allow the user
        // to request to join companies where they have been previously rejected.
        associatedCompanies: this.rows.filter(row => row.Status === 'pending' || row.Status === 'Joined' || row.Status === 'approved')
        //associatedCompanies: this.rows
      });
      modal.present();
    }
  }
}
