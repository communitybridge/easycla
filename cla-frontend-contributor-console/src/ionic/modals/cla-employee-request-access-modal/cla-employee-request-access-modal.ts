// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { ChangeDetectorRef, Component } from '@angular/core';
import { AlertController, IonicPage, ModalController, NavController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company/contact'
})
@Component({
  selector: 'cla-employee-request-access-modal',
  templateUrl: 'cla-employee-request-access-modal.html',
})
export class ClaEmployeeRequestAccessModal {
  project: any;
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  company: any;
  authenticated: boolean;
  cclaSignature: any;
  managers: any;

  userEmails: Array<string>;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  loading: any;
  showManagerSelectOption: boolean;
  showManagerEnterOption: boolean;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private formBuilder: FormBuilder,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.loading = {
      signatures: true,
    };
    this.project = {};
    this.company = {};

    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.authenticated = navParams.get('authenticated');
    this.form = formBuilder.group({
      email: [''], // Validators.compose([Validators.required, EmailValidator.isValid])
      manager: [''], // Validators.compose([Validators.required])
      message: [''], // Validators.compose([Validators.required])
      adminname: [''],
      managerOptions: [''],
      adminemail: [''], // Validators.compose([Validators.required, EmailValidator.isValid])
    });
    this.managers = [];
}

saveManagerOption() {
  const option = this.form.value.managerOptions;
  if (option === 'select manager') {
    this.showManagerSelectOption = true;
    this.showManagerEnterOption = false;
  }
  else if (option === 'enter manager') {
    this.showManagerSelectOption = false;
    this.showManagerEnterOption = true;
  }
}

getDefaults() {
  this.userEmails = [];
}

ngOnInit() {
  this.getUser(this.userId, this.authenticated).subscribe(user => {
    if (user) {
      this.userEmails = user.user_emails || [];
      if (user.lf_email && this.userEmails.indexOf(user.lf_email) == -1) {
        this.userEmails.push(user.lf_email)
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
  this.claService.getProject(projectId).subscribe(response => {
    this.project = response;
  });
}

getCompany(companyId: string) {
  this.claService.getCompany(companyId).subscribe(response => {
    this.company = response;
  });
}

getProjectSignatures(projectId: string, companyId: string) {
  // Get CCLA Company Signatures - should just be one
  this.loading.signatures = true;
  this.claService.getCompanyProjectSignatures(companyId, projectId)
    .subscribe(response => {
      this.loading.signatures = false;
      console.log('Signatures for project: ' + projectId + ' for company: ' + companyId);
      console.log(response);
      if (response.signatures) {
        let cclaSignatures = response.signatures.filter(sig => sig.signatureType === 'ccla');
        console.log('CCLA Signatures for project: ' + cclaSignatures.length);
        if (cclaSignatures.length) {
          console.log('CCLA Signatures for project id: ' + projectId + ' and company id: ' + companyId);
          console.log(cclaSignatures);
          this.cclaSignature = cclaSignatures[0];
          console.log(this.cclaSignature);
          console.log(this.cclaSignature.signatureACL);
          if (this.cclaSignature.signatureACL != null) {
            for (let manager of this.cclaSignature.signatureACL) {
              this.managers.push({
                userID: manager.userID,
                username: manager.username,
                lfEmail: manager.lfEmail,
              });
            }
          }
        }
      }
    },
      exception => {
        this.loading.signatures = false;
        console.log("Exception while calling: getCompanyProjectSignatures() for company ID: " +
          companyId + ' and project ID: ' + projectId);
        console.log(exception);
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
  
  let message = {
    user_email: this.form.value.email,
    message: this.form.value.message,
    project_id: this.projectId,
  };
  
  // this.claService.postUserMessageToCompanyManager(this.userId, this.companyId, message).subscribe(response => {
  //   this.emailSent();
  // });
}

emailSent() {
  let message = this.authenticated ?
    'Thank you for contacting your company\'s administrators. Once the CLA is signed and you are authorized, please navigate to the Agreements tab in the Gerrit Settings page and restart the CLA signing process' :
    'Thank you for contacting your company\'s administrators. Once the CLA is signed and you are authorized, you will have to complete the CLA process from your existing pull request.'
  let alert = this.alertCtrl.create({
    title: 'E-Mail Successfully Sent!',
    subTitle: message,
    buttons: ['Dismiss']
  });
  alert.onDidDismiss(() => this.dismiss());
  alert.present();
}
}
