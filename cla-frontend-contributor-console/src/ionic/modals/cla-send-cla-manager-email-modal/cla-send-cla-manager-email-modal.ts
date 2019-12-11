// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company/contact'
})
@Component({
  selector: 'cla-send-cla-manager-email-modal',
  templateUrl: 'cla-send-cla-manager-email-modal.html'
})
export class ClaSendClaManagerEmailModal {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  authenticated: boolean;

  company: any;

  userEmails: Array<string>;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.authenticated = navParams.get('authenticated');
    this.form = formBuilder.group({
      email: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      message: [''] // Validators.compose([Validators.required])
    });
  }

  getDefaults() {
    this.userEmails = [];
    this.company = {
      company_name: ''
    };
  }

  ngOnInit() {
    this.getUser();
    this.getCompany();
  }

  getUser() {
    if (this.authenticated) {
      // Gerrit Users
      this.claService.getUserWithAuthToken(this.userId).subscribe(user => {
        if (user) {
          this.userEmails = user.user_emails || [];
          if (user.lf_email && this.userEmails.indexOf(user.lf_email) == -1) {
            this.userEmails.push(user.lf_email);
          }
        } else {
          console.log('Unable to retrieve user.');
        }
      });
    } else {
      // Github Users
      this.claService.getUser(this.userId).subscribe(user => {
        if (user) {
          this.userEmails = user.user_emails || [];
        } else {
          console.log('Unable to retrieve user.');
        }
      });
    }
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      this.company = response;
    });
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let data = {
      user_email: this.form.value.email,
      company_id: this.companyId,
      project_id: this.projectId
    };
    this.claService.postUserCCLARequestToManager(this.userId, data).subscribe(response => {
      this.emailSent();
    });
  }

  emailSent() {
    let alert = this.alertCtrl.create({
      title: 'E-Mail Successfully Sent!',
      subTitle:
        'Thank you for contacting your CLA Manager. Once you are authorized, you will have to complete the CLA process from your existing pull request.',
      buttons: ['Dismiss']
    });
    alert.onDidDismiss(() => this.dismiss());
    alert.present();
  }
}
