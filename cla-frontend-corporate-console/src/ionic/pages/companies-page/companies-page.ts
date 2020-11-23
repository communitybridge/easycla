// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component, OnInit} from '@angular/core';
import {IonicPage, ModalController, NavController} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';
import {Restricted} from '../../decorators/restricted';
import {ClaCompanyWithInvitesModel} from "../../models/cla-company-with-invites";

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
export class CompaniesPage implements OnInit {
  loading: any;
  companies: ClaCompanyWithInvitesModel[] = [];
  pendingRejectedCompanies: ClaCompanyWithInvitesModel[] = [];
  userId: string;
  userEmail: string;
  userName: string;
  companyId: string;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  expanded: boolean = true;

  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companies: true
    };
    this.userId = localStorage.getItem('userid');
    if (this.userId != null) {
      this.userId = this.userId.trim();
    }
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
        if (exception.status === 404) {
          // Create the user if it does not exist
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
      (response) => {
        this.loading.companies = false;
        if (response['companies-with-invites']) {
          let companyModels: ClaCompanyWithInvitesModel[] = this.mapCompanies(response['companies-with-invites']);
          for (let companyModel of companyModels) {
            if (this.userInCompanyACL(companyModel)) {
              this.companies.push(companyModel);
            } else {
              this.pendingRejectedCompanies.push(companyModel);
            }
          }
        } else {
          this.companies = [];
          this.pendingRejectedCompanies = [];
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

  mapCompanies(companyResponseModels: any): ClaCompanyWithInvitesModel[] {
    let companyModels: ClaCompanyWithInvitesModel[] = [];
    for (let companyResponseModel of companyResponseModels) {
      let companyModel = new ClaCompanyWithInvitesModel(
        companyResponseModel.companyID,
        companyResponseModel.companyName,
        companyResponseModel.companyACL,
        companyResponseModel.status,
        companyResponseModel.created,
        companyResponseModel.updated,
      );
      companyModels.push(companyModel);
    }

    // Sort and remove duplicates
    companyModels = this.sortData(companyModels);
    companyModels = this.removeDuplicates(companyModels);

    return companyModels;
  }

  sortData(companies: ClaCompanyWithInvitesModel[]) {
    let joinedCompanies = companies.filter(company => (company.status !== 'pending' && company.status !== 'rejected'))
    let requestCompanies = companies.filter(company => company.status === 'pending')
    joinedCompanies = joinedCompanies.sort((a, b) => {
      return a.name.toLowerCase().localeCompare(b.name.toLowerCase());
    });
    requestCompanies = requestCompanies.sort((a, b) => {
      return a.name.toLowerCase().localeCompare(b.name.toLowerCase());
    });

    return joinedCompanies.concat(requestCompanies);
  }

  removeDuplicates(companies: ClaCompanyWithInvitesModel[]) {
    interface SeenMapType {
      [key: string]: boolean;
    }

    let seenCompanies: SeenMapType = {};
    return companies.filter(company => {
      if (seenCompanies[company.id] == null) {
        seenCompanies[company.id] = true;
        return true; // unique, pass filter
      }
      return false; // duplicate, fail filter
    });
  }

  openSelectCompany() {
    if (!this.loading.companies) {
      let modal = this.modalCtrl.create('AddCompanyModal', {
        // For this modal, we share the list of companies where the user is either pending acceptance, approved/joined.
        // From this, the modal dialog will filter the list so that these companies are not shown. We allow the user
        // to request to join companies where they have been previously rejected.
        associatedCompanies: this.companies.filter(company => company.status === 'pending' || company.status === 'Joined' || company.status === 'approved')
      });

      modal.present();
    }
  }

  onClickToggle(hasExpanded) {
    this.expanded = hasExpanded;
  }

  userInCompanyACL(company: ClaCompanyWithInvitesModel) {
    if (company.acl) {
      return company.acl.indexOf(this.userId) > -1;
    }
    return false;
  }
}
