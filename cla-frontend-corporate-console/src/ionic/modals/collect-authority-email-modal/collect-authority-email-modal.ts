// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/collect-authority-email'
})
@Component({
  selector: 'collect-authority-email-modal',
  templateUrl: 'collect-authority-email-modal.html'
})
export class CollectAuthorityEmailModal {
  projectId: string;
  companyId: string;
  signingType: string;
  projectName: string;
  companyName: string;
  form: FormGroup;
  submitAttempt: boolean = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.projectName = navParams.get('projectName');
    this.companyName = navParams.get('companyName');
    this.projectId = navParams.get('projectId');
    this.companyId = navParams.get('companyId');
    this.signingType = navParams.get('signingType');
    this.form = formBuilder.group({
      authorityemail: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
      authorityname: ['']
    });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  emailSent() {
    let alert = this.alertCtrl.create({
      title: 'E-Mail Sent!',
      subTitle: 'An E-Mail has been sent. Please wait for your CLA Signatory to review and sign the CLA.',
      buttons: [
        {
          text: 'Dismiss',
          role: 'dismiss',
          handler: () => {
            this.navCtrl.pop();
          }
        }
      ]
    });
    alert.present();
  }

  submit() {
    this.submitAttempt = true;
    if (!this.form.valid) {
      return;
    }
    let emailRequest = {
      project_id: this.projectId,
      company_id: this.companyId,
      send_as_email: true,
      authority_name: this.form.value.authorityname,
      authority_email: this.form.value.authorityemail
    };

    this.claService.postCorporateSignatureRequest(emailRequest).subscribe((response) => {
      this.emailSent();
      this.dismiss();
    });
  }
}
