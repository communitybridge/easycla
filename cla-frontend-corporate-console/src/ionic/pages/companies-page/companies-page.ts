// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';
import { ColumnMode, SelectionType, SortType } from '@swimlane/ngx-datatable';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';

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

  manager: string;
  columns: any[];
  rows: any[];
  formErrors: any[];
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  formSuccessfullySubmitted: boolean = false;
  claManagerApproved: boolean = false;

  ColumnMode = ColumnMode;
  SelectionType = SelectionType;
  SortType = SortType;

  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private formBuilder: FormBuilder,
    private rolesService: RolesService // for @Restricted
  ) {
    this.form = formBuilder.group({
      project_name: [''],
      compnay_name: ['', Validators.compose([Validators.required])],
      full_name: ['', Validators.compose([Validators.required])],
      lfid: [''],
      email_address: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      cla_admin: ['', Validators.compose([Validators.required])]
    });
    this.formErrors = [];
    this.getDefaults();
  }

  openCLAOnboardingForm() {
    // this.navCtrl.push('ClaManagerOnboardingPage');
  }

  approveCLAManager() {
    this.claManagerApproved = true;
  }

  getDefaults() {
    this.loading = {
      companies: true
    };
    this.userId = localStorage.getItem('userid');
    this.userEmail = localStorage.getItem('user_email');
    this.userName = localStorage.getItem('user_name');
    this.setUserDetails();
    this.companies = [];
    this.columns = [{ prop: 'CompanyName' }, { prop: 'Status' }, { prop: 'Action' }, { prop: 'CompanyID' }, { prop: 'ProjectName' }];
  }

  ngOnInit() {
    this.getCompanies();
    this.getSignedCLAs();
  }

  setUserDetails() {
    this.form.controls['lfid'].setValue(this.userId);
    this.form.controls['email_address'].setValue(this.userEmail);
    this.form.controls['full_name'].setValue(this.userName);
  }

  openCompanyModal() {
    let modal = this.modalCtrl.create('AddCompanyModal', {});
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getCompanies();
    });
    modal.present();
  }

  getSignedCLAs() {}

  getPendingInvites() {}

  /**
   * Get the list of companies for this user - if any
   */
  getCompanies() {
    this.loading.companies = true;
    this.claService.getUserByUserName(this.userId).subscribe(
      (response) => {
        //console.log(response);
        // We need the user's unique ID - grab it from the first record
        this.getCompaniesByUserManagerWithInvites(response.userID);
      },
      (exception) => {
        // Typically get this if the user is not found in our database
        this.loading.companies = false;
        this.rows = [];

        console.log('Received exception:');
        console.log(exception);
        if (exception.status === 404) {
          // Create the user if it does't exist
          console.log('Creating user record...');
          const user = {
            lfUsername: this.userId,
            username: this.userName,
            lfEmail: this.userEmail
          };
          this.claService.createUserV3(user).subscribe(
            (response) => {
              console.log('Success creating user record: ');
              console.log(response);
            },
            (exception) => {
              console.log('Error creating user record: ');
              console.log(exception);
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
          // console.log(this.rows, 'this is rows')
        } else {
          this.rows = [];
        }
      },
      (exception) => {
        this.loading.companies = false;
        console.log('Exception while calling: getCompaniesByUserManagerWithInvites() for userId: ' + userId);
        console.log(exception);
      }
    );
  }

  viewCompany(companyId) {
    this.navCtrl.setRoot('CompanyPage', {
      companyId: companyId
    });
  }

  onSelect(event) {
    let company = event.selected[0];
    if (company.Status === 'Joined') {
      this.viewCompany(company.CompanyID);
    }
  }

  mapCompanies(companies) {
    let rows = [];
    let action;
    for (let company of companies) {
      if (company.status === 'Pending Approval') {
        action = '';
      }
      rows.push({
        CompanyID: company.companyID,
        CompanyName: company.companyName,
        Status: company.status,
        ProjectName: ''
      });
    }
    return rows;
  }
}
