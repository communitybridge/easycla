// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input } from '@angular/core';
import { AlertController, IonicPage, ViewController, NavParams } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators, FormControl } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { ClaCompanyModel } from '../../models/cla-company';

@IonicPage({
  segment: 'edit-company-modal'
})
@Component({
  selector: 'edit-company-modal',
  templateUrl: 'edit-company-modal.html'
})
export class EditCompanyModal {
  @Input() company: any;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService,
    public alertCtrl: AlertController,
    private navParams: NavParams
  ) {
    this.company = this.navParams.get('company');
  }

  get companyName(): FormControl {
    return <FormControl>this.form.get('companyName');
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      companyName: [this.company.company_name, Validators.compose([Validators.required])]
    });
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (this.form.valid) {
      this.updateCompany();
    } else {
      this.currentlySubmitting = false;
    }
  }

  updateCompany() {
    let company = {
      company_id: this.company.company_id,
      company_name: this.companyName.value
    };
    this.claService.putCompany(company).subscribe(
      () => {
        this.currentlySubmitting = false;
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

  dismiss() {
    this.viewCtrl.dismiss();
  }

  companyExistAlert(company_id) {
    let alert = this.alertCtrl.create({
      title: 'Company ' + this.companyName.value + ' already exists',
      message: 'The company you tried to create already exists in the CLA system. Would you like to request access?',
      buttons: [
        {
          text: 'Request',
          handler: () => {
            const userId = localStorage.getItem('userid');
            const userEmail = localStorage.getItem('user_email');
            const userName = localStorage.getItem('user_name');
            this.claService
              .sendInviteRequestEmail(company_id, userId, userEmail, userName)
              .subscribe(() => this.dismiss());
          }
        },
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('No clicked');
          }
        }
      ]
    });
    alert.present();
  }
}
