// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { AlertController, IonicPage, ModalController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';
import { Content } from 'ionic-angular';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/invite-company-admin'
})
@Component({
  selector: 'cla-company-admin-send-email-modal',
  templateUrl: 'cla-company-admin-send-email-modal.html'
})
export class ClaCompanyAdminSendEmailModal {
  projectId: string;
  userId: string;
  authenticated: boolean; // true if coming from gerrit/corporate
  userEmails: Array<string>;
  form: FormGroup;
  serverError: string = '';
  isSendClicked = false;
  @ViewChild('pageTop') pageTop: Content;

  constructor(
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.userEmails = [];
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.authenticated = navParams.get('authenticated');
    this.form = formBuilder.group({
      company_name: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      contributor_name: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      contributor_email: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      cla_manager_name: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      cla_manager_email: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
    });
  }

  ngOnInit() {
    this.getUser();
  }

  getUser() {
    if (this.authenticated) {
      // Gerrit Users
      this.claService.getUserWithAuthToken(this.userId).subscribe((user) => {
        if (user) {
          this.userEmails = user.user_emails || [];
          if (user.lf_email && this.userEmails.indexOf(user.lf_email) == -1) {
            this.userEmails.push(user.lf_email);
          }
        } else {
          console.warn('Unable to retrieve user.');
        }
      });
    } else {
      // Github Users
      this.claService.getUser(this.userId).subscribe((user) => {
        if (user) {
          this.userEmails = user.user_emails || [];
        } else {
          console.warn('Unable to retrieve user.');
        }
      });
    }
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  emailSent() {
    let alert = this.alertCtrl.create({
      title: 'E-Mail Sent!',
      subTitle: 'An E-Mail has been sent. Please wait for your CLA Manager to add you to your company approved list.',
      buttons: ['Dismiss']
    });
    alert.present();
  }

  submit() {
    this.isSendClicked = true;
    if (this.form.valid) {
      this.claService.getProject(this.projectId).subscribe((response) => {
        this.sendRequest(response);
      });
    }
  }

  sendRequest(project) {
    this.serverError = '';
    let data = {
      contributor_name: this.form.value.contributor_name,
      contributor_email: this.form.value.contributor_name,
      cla_manager_name: this.form.value.cla_manager_name,
      cla_manager_email: this.form.value.cla_manager_email,
      project_name: project.project_name,
      company_name: this.form.value.company_name,
    };
    this.claService.postEmailToCompanyAdmin(this.userId, data).subscribe(
      (response) => {
        this.isSendClicked = false;
        this.emailSent();
        this.dismiss();
      },
      (exception) => {
        this.isSendClicked = false;
        const errorObj = JSON.parse(exception._body);
        if (errorObj) {
          this.serverError = errorObj.errors.contributor_email;
          this.pageTop.scrollToTop();
        }
      }
    );
  }
}
