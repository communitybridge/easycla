// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {AlertController, IonicPage, ModalController, NavParams, ViewController} from 'ionic-angular';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {EmailValidator} from '../../validators/email';
import {ClaService} from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company/contact'
})
@Component({
  selector: 'cla-employee-request-access-modal',
  templateUrl: 'cla-employee-request-access-modal.html'
})
export class ClaEmployeeRequestAccessModal {
  project: any;
  projectId: string;
  userId: string;
  companyId: string;
  company: any;
  authenticated: boolean;
  cclaSignature: any;
  managers: any;
  formErrors: any[];

  userEmails: Array<string>;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  loading: any;
  showManagerSelectOption: boolean;
  showManagerEnterOption: boolean;

  constructor(
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.getDefaults();
    this.loading = true;
    this.project = {};
    this.company = {};

    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.authenticated = navParams.get('authenticated');

    this.form = formBuilder.group({
      user_email: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      user_name: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      message: ['', Validators.compose([Validators.required])],
      recipient_name: [''],
      recipient_email: [''],
      manager: [''],
      managerOptions: ['', Validators.compose([Validators.required])]
    });
    this.managers = [];
    this.formErrors = [];
  }

  saveManagerOption() {
    const option = this.form.value.managerOptions;
    if (option === 'select manager') {
      this.showManagerSelectOption = true;
      this.showManagerEnterOption = false;
      this.resetFormValues('recipient_name');
      this.resetFormValues('recipient_email');
      this.form.controls['recipient_name'].clearValidators();
      this.form.controls['recipient_email'].clearValidators();
    } else if (option === 'enter manager') {
      this.showManagerSelectOption = false;
      this.showManagerEnterOption = true;
      if (this.managers.length > 1) {
        this.resetFormValues('manager');
      }
      this.form.controls['recipient_name'].setValidators(Validators.compose([Validators.required]));
      this.form.controls['recipient_email'].setValidators(Validators.compose([Validators.required, EmailValidator.isValid]));
    }
    this.form.controls['recipient_name'].updateValueAndValidity();
    this.form.controls['recipient_email'].updateValueAndValidity();
  }

  resetFormValues(value) {
    return this.form.controls[value].reset();
  }

  getCLAManagerDetails(managerId) {
    return this.managers.find((manager) => {
      return (manager.userID === managerId);
    });
  }

  getDefaults() {
    this.userEmails = [];
  }

  ngOnInit() {
    this.getUser(this.userId, this.authenticated).subscribe((user) => {
      if (user) {
        if (user.user_emails.length === 1) {
          this.form.controls['user_email'].setValue(user.user_emails[0]);
        }
        this.userEmails = user.user_emails || [];
        if (user.lf_email && this.userEmails.indexOf(user.lf_email) == -1) {
          this.userEmails.push(user.lf_email);
        }
      }
    });
    this.getProject(this.projectId);
    this.getCompany(this.companyId);
    this.getProjectSignatures(this.projectId, this.companyId);
  }

  getUser(userId: string, authenticated: boolean) {
    if (authenticated) {
      // Gerrit Users
      return this.claService.getUserWithAuthToken(userId);
    } else {
      // Github Users
      return this.claService.getUser(userId);
    }
  }

  getProject(projectId: string) {
    this.claService.getProject(projectId).subscribe((response) => {
      this.project = response;
    });
  }

  getCompany(companyId: string) {
    this.claService.getCompany(companyId).subscribe((response) => {
      this.company = response;
    });
  }

  insertAndSortManagersList(manager) {
    this.managers.push(manager);
    this.managers.sort((first, second) => {
      return first.username.toLowerCase() - second.username.toLowerCase();
    });
  }

  getProjectSignatures(projectId: string, companyId: string) {
    // Get CCLA Company Signatures - should just be one
    this.loading = true;
    this.claService.getCompanyProjectSignatures(companyId, projectId).subscribe(
      (response) => {
        this.loading = false;
        if (response.signatures) {
          let cclaSignatures = response.signatures.filter((sig) => sig.signatureType === 'ccla');
          if (cclaSignatures.length) {
            this.cclaSignature = cclaSignatures[0];
            if (this.cclaSignature.signatureACL != null) {
              if (this.cclaSignature.signatureACL.length === 1) {
                this.form.controls['manager'].setValue(this.cclaSignature.signatureACL[0].userID);
              }
              for (let manager of this.cclaSignature.signatureACL) {
                this.insertAndSortManagersList({
                  userID: manager.userID,
                  username: manager.username,
                  lfEmail: manager.lfEmail
                });
              }
            }
          }
        }
      },
      (exception) => {
        this.loading = false;
      }
    );
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    this.formErrors = [];

    if (!this.form.valid) {
      this.getFormValidationErrors();
      this.currentlySubmitting = false;
      return;
    }

    let managerEmail = '';
    let managerUsername = '';

    // 'select manager' or 'enter manager'
    if (this.form.value.manager && this.form.value.managerOptions === 'select manager') {
      managerEmail = this.getCLAManagerDetails(this.form.value.manager).lfEmail;
      managerUsername = this.getCLAManagerDetails(this.form.value.manager).username;
    } else {
      managerEmail = this.form.value.recipient_email;
      managerUsername = this.form.value.recipient_name;
    }

    let data = {
      contributorId: this.userId,
      contributorName: this.form.value.user_name,
      contributorEmail: this.form.value.user_email,
      message: this.form.value.message,
      recipientName: managerUsername,
      recipientEmail: managerEmail,
    };

    this.claService.requestToBeOnCompanyApprovedList(this.userId, this.companyId, this.projectId, data)
      .subscribe((response) => {
        this.loading = true;
        this.emailSent();
      }, (error) => {
        this.loading = true;
        this.emailSentError(error);
      });
  }

  emailSent() {
    this.loading = false;
    let message = this.authenticated
      ? "Thank you for contacting your company's administrators. Once the CLA is signed and you are authorized, please navigate to the Agreements tab in the Gerrit Settings page and restart the CLA signing process"
      : "Thank you for contacting your company's administrators. Once the CLA is signed and you are authorized, you will have to complete the CLA process from your existing pull request.";
    let alert = this.alertCtrl.create({
      title: 'E-Mail Successfully Sent!',
      subTitle: message,
      buttons: ['Dismiss']
    });
    alert.onDidDismiss(() => this.dismiss());
    alert.present();
  }

  emailSentError(error) {
    this.loading = false;
    let message = `The request already exists for you. Please ask the CLA Manager to log into the EasyCLA Corporate Console and authorize you using one of the available methods.`;
    let alert = this.alertCtrl.create({
      title: 'Problem Sending Request',
      subTitle: message,
      buttons: ['Dismiss']
    });
    alert.onDidDismiss(() => this.dismiss());
    alert.present();
  }

  getFormValidationErrors() {
    let message;
    Object.keys(this.form.controls).forEach((key) => {
      const controlErrors = this.form.get(key).errors;
      if (controlErrors != null) {
        Object.keys(controlErrors).forEach((keyError) => {
          switch (key) {
            case 'managerOptions':
              message = `*Selecting an Option for Entering a CLA Manager is ${keyError}`;
              break;
            case 'user_name':
              message = `*User Name Field is ${keyError}`;
              break;
            case 'user_email':
              message = `*Email Authorize Field is ${keyError}`;
              break;
            case 'recipient_email':
              message = `*Receipent Email Field is ${keyError}`;
              break;
            case 'message':
              message = `*Message Field is ${keyError}`;
              break;

            default:
              message = `Check Fields for errors`;
          }
          this.formErrors.push({
            message
          });
        });
      }
    });
  }

  trimCharacter(text, length) {
    return text.length > length ? text.substring(0, length) + '...' : text;
  }
}
