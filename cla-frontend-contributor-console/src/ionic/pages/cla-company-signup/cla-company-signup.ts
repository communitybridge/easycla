// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from '../../validators/email';

// @IonicPage({
//   segment: 'cla/project/:projectId/company-signup'
// })
@Component({
  selector: 'cla-company-signup',
  templateUrl: 'cla-company-signup.html'
})
export class ClaCompanySignupPage {
  projectId: string;
  repositoryId: string;

  project: any;
  company: any;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private formBuilder: FormBuilder
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');

    this.form = formBuilder.group({
      company: ['', Validators.compose([Validators.required])],
      email: ['', Validators.compose([Validators.required, EmailValidator.isValid])]
    });
  }

  getDefaults() {}

  ngOnInit() {
    this.getProject();
  }

  getProject() {
    this.project = {
      id: '0000000001',
      name: 'Project Name',
      logoRef: 'https://dummyimage.com/225x102/d8d8d8/242424.png&text=Project+Logo'
    };
  }

  getCompany() {
    this.company = {
      name: 'Company Name',
      id: '0000000001'
    };
  }

  submit() {
    // Need to Post a request to cla to generate a new docusign document
    // should return a url to that contract
    // We'll need to open a new page and load the link into the page
  }
}
