// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { AlertController, IonicPage, ModalController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';
import { Content } from 'ionic-angular';
import { generalConstants } from '../../constants/general';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/invite-company-admin'
})
@Component({
  selector: 'cla-company-admin-send-email-modal',
  templateUrl: 'cla-company-admin-send-email-modal.html'
})
export class ClaCompanyAdminSendEmailModal {
  projectId: string;
  companyId: string;
  companyName: string;
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
    // May be empty
    this.companyId = navParams.get('companyId');
    // May be empty
    this.companyName = navParams.get('companyName');
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
    this.getUserEmails();
  }

  getUserEmails() {
    const user = JSON.parse(localStorage.getItem(generalConstants.USER_MODEL));
    if (user) {
      this.userEmails = user.user_emails || [];
      if (user.lf_email && this.userEmails.indexOf(user.lf_email) == -1) {
        this.userEmails.push(user.lf_email);
      }
    } else {
      console.warn('Unable to retrieve user.');
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
        // Instead of creating a company we need to send email to CLA Manager.
        this.sendRequest(response);
      });
    }
  }

  sendRequest(project) {
    this.serverError = '';
    let data = {
      contributorName: this.form.value.contributor_name,
      contributorEmail: this.form.value.contributor_email,
      claManagerName: this.form.value.cla_manager_name,
      claManagerEmail: this.form.value.cla_manager_email,
      projectName: project.project_name,
      companyName: this.form.value.company_name,
      version: 'v1'
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
          this.serverError = errorObj.Message;
          this.pageTop.scrollToTop();
        }
      }
    );
  }
}
