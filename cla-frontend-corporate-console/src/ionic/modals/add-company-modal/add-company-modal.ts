// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { AlertController, IonicPage, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { ClaCompanyModel } from '../../models/cla-company';
import { AuthService } from '../../services/auth.service';

@IonicPage({
  segment: 'add-company-modal'
})
@Component({
  selector: 'add-company-modal',
  templateUrl: 'add-company-modal.html'
})
export class AddCompanyModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  company: ClaCompanyModel;
  companyName: string;
  userEmail: string;
  userId: string;
  userName: string;
  companies: any[];
  filteredCompanies: any[];
  companySet: boolean = false;
  addNewCompany: boolean = false;
  existingCompanyId: string;
  mode: string = 'add';
  loading: any;
  searching: boolean;
  activateButtons: boolean;
  join: boolean;
  add: boolean;
  associatedCompanies: any[] = [];

  constructor(
    private navParams: NavParams,
    private viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    private claService: ClaService,
    private authService: AuthService,
    private alertCtrl: AlertController
  ) {
    this.associatedCompanies = this.navParams.get('associatedCompanies');
    this.getDefaults();
  }

  getDefaults() {
    this.searching = false;
    this.userName = localStorage.getItem('user_name');
    this.userId = localStorage.getItem('userid');
    this.company = this.navParams.get('company');
    this.mode = this.navParams.get('mode') || 'add';
    this.companies = [];
    this.filteredCompanies = [];
    this.loading = {
      submit: false,
      companies: true
    };
    this.addNewCompany = true;
    this.add = true;
    this.join = false;
    this.activateButtons = true;

    this.form = this.formBuilder.group({
      companyName: [this.companyName, Validators.compose([Validators.required, Validators.maxLength(60)])],
    });
  }

  ngOnInit() {
    this.getAllCompanies();
  }

  ionViewDidEnter() {
    this.updateUserInfoBasedLFID();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    this.addCompany();
  }

  addCompany() {
    this.loading.submit = true;
    let company = {
      company_name: this.companyName,
      company_manager_user_email: this.userEmail,
      company_manager_user_name: this.userName,
      company_manager_id: this.userId
    };
    this.claService.postCompany(company).subscribe(
      (response) => {
        this.currentlySubmitting = false;
        window.location.reload(true);
        this.dismiss();
      },
      (err: any) => {
        if (err.status === 409) {
          let errJSON = err.json();
          this.companyExistAlert(errJSON.company_id);
        }
        this.currentlySubmitting = false;
      }
    );
  }

  sendCompanyNotification(reload) {
    this.loading.submit = false;
    let alert = this.alertCtrl.create({
      title: 'Notification Sent!',
      subTitle: `A Notification has been sent to the Company Administrators for ${this.companyName}`,
      buttons: [
        {
          text: 'Ok',
          role: 'dismiss',
        }
      ]
    });
    alert.onDidDismiss(() => {
      if (reload) {
        window.location.reload(true);
      }
    });
    alert.present();
  }

  joinCompany() {
    this.loading.submit = true;
    this.sendInvitationForCompany(this.existingCompanyId);
  }

  dismiss() {
    this.viewCtrl.dismiss(this.existingCompanyId);
  }

  companyExistAlert(companyId) {
    let alert = this.alertCtrl.create({
      title: 'Company ' + this.companyName + ' already exists',
      message: 'The company you tried to create already exists in the CLA system. Would you like to request access?',
      buttons: [
        {
          text: 'Request',
          handler: () => {
            this.sendInvitationForCompany(companyId);
          }
        },
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            this.loading.submit = false;
          }
        }
      ]
    });
    alert.present();
  }

  sendInvitationForCompany(companyId) {
    const userId = localStorage.getItem('userid');
    const userEmail = localStorage.getItem('user_email');
    const userName = localStorage.getItem('user_name');
    this.claService
      .sendInviteRequestEmail(companyId, userId, userEmail, userName)
      .subscribe(() => {
        this.sendCompanyNotification(true);
        this.dismiss();
      });
  }

  getAllCompanies() {
    if (!this.companies) {
      this.loading.companies = true;
    }
    this.claService.getAllV3Companies().subscribe((response) => {
      this.loading.companies = false;
      this.associatedCompanies = this.associatedCompanies.map(aCompany => aCompany.CompanyID);
      this.companies = response.companies.filter(company => this.associatedCompanies.indexOf(company.companyID) < 0);
    });
  }

  findCompany(event) {
    this.filteredCompanies = [];
    if (!this.companies) {
      this.searching = true;
    }

    if (!this.companySet) {
      this.join = false;
      this.add = true;
    } else {
      this.join = true;
      this.add = false;
    }

    let companyName = event.value;
    if (companyName.length > 0 && this.companies) {
      this.activateButtons = false;
      this.searching = false;
      this.companySet = false;
      this.filteredCompanies = this.companies
        .map((company) => {
          let formattedCompany;
          if (company.companyName.toLowerCase().includes(companyName.toLowerCase())) {
            formattedCompany = company.companyName.replace(
              new RegExp(companyName, 'gi'),
              (match) => '<span class="highlightText">' + match + '</span>'
            );
          }
          company.filteredCompany = formattedCompany;
          return company;
        })
        .filter((company) => company.filteredCompany);
    } else {
      this.activateButtons = true;
    }

    if (companyName.length >= 2) {
      this.addNewCompany = false;
    }
  }

  setCompanyName(company) {
    this.companySet = true;
    this.companyName = company.companyName;
    this.existingCompanyId = company.companyID;
  }

  private updateUserInfoBasedLFID() {
    if (this.authService.isAuthenticated()) {
      this.authService
        .getIdToken()
        .then((token) => {
          return this.authService.parseIdToken(token);
        })
        .then((tokenParsed) => {
          if (tokenParsed && tokenParsed['email']) {
            this.userEmail = tokenParsed['email'];
          }
          if (tokenParsed && tokenParsed['name']) {
            this.userName = tokenParsed['name'];
          }
        })
        .catch((error) => {
          console.warn(JSON.stringify(error));
          return;
        });
    }
    return;
  }
}
