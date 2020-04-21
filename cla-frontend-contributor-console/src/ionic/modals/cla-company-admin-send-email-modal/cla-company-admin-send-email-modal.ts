// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';

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
  submitAttempt: boolean = false;

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
      useremail: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      adminemail: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      adminname: [''],
      username: ['']
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
    this.submitAttempt = true;
    if (!this.form.valid) {
      return;
    }

    this.claService.getProject(this.projectId).subscribe((project) => {
      // TODO - Add company_name to the data payload
      let data = {
        user_email: this.form.value.useremail,
        cla_manager_name: this.form.value.adminname,
        cla_manager_email: this.form.value.adminemail,
        project_name: project.project_name,
        company_name: "",
        user_name: this.form.value.username,
      };
      this.claService.postEmailToCompanyAdmin(this.userId, data).subscribe((response) => {
        this.emailSent();
        this.dismiss();
      });
    });
  }
}
