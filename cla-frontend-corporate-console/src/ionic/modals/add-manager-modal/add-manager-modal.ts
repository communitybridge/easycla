// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { generalConstants } from '../../constant/general';
import { EmailValidator } from '../../validators/email';
import { LFIDValidator } from '../../validators/lfid';
import { UserNameValidator } from '../../validators/user-name';

@IonicPage({
  segment: 'add-manager-modal'
})
@Component({
  selector: 'add-manager-modal',
  templateUrl: 'add-manager-modal.html'
})
export class AddManagerModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  projectId: string;
  companyId: string;
  signatureId: string;
  managerName: string;
  managerEmail: string;
  managerLFID: string;
  errorMsg: string;
  linuxFoundationIdentityURL: string = generalConstants.linuxFoundationIdentityURL;

  constructor(
    private viewCtrl: ViewController,
    private navParams: NavParams,
    private formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.projectId = this.navParams.get('projectId');
    this.companyId = this.navParams.get('companyId');
    this.signatureId = this.navParams.get('signatureId');
    this.form = this.formBuilder.group({
      managerName: [this.managerName, Validators.compose([Validators.required, Validators.minLength(3), UserNameValidator.isValid])],
      managerLFID: [this.managerLFID, Validators.compose([Validators.required, LFIDValidator.isValid])],
      managerEmail: [this.managerEmail, Validators.compose([Validators.required, EmailValidator.isValid])],
    });
  }

  submit() {
    if (/[~`!@#$%\^&*()+=\-\[\]\\';,/{}|\\":<>\?]/g.test(this.form.value.managerLFID)) {
      this.errorMsg = 'Special characters are not allowed';
      return;
    }
    this.submitAttempt = true;
    this.addManager();
  }

  addManager() {
    this.claService.addCLAManager(this.companyId, this.projectId, this.form.getRawValue()).subscribe((result) => {
      if (result.errors) {
        this.errorMsg = result.errors[0];
      } else {
        this.dismiss(true);
      }
    });
  }

  dismiss(data = false) {
    this.viewCtrl.dismiss(data);
  }

  clearError(event) {
    this.errorMsg = '';
  }
}
