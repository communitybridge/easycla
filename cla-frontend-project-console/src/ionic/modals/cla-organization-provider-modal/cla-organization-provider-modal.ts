// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, Events } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { Http } from '@angular/http';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla-organization-provider-modal'
})
@Component({
  selector: 'cla-organization-provider-modal',
  templateUrl: 'cla-organization-provider-modal.html'
})
export class ClaOrganizationProviderModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  responseErrors: string[] = [];
  claProjectId: any;
  showErrorMsg: boolean;
  loading: boolean;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public http: Http,
    public claService: ClaService,
    public events: Events
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.form = formBuilder.group({
      // provider: ['', Validators.required],
      orgName: ['', Validators.compose([Validators.required]) /*, this.urlCheck.bind(this)*/]
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.checkGitOrganisationName();
  }

  checkGitOrganisationName() {
    this.loading = true;
    this.showErrorMsg = false;
    this.claService.testGitHubOrganization(this.form.value.orgName).subscribe((res: any)=> {
      this.loading = false;
      if(res.status ===  200) {
        this.postClaGithubOrganization();
      }
    }, (err: any) => { 
      this.loading = false;
      if(!err.ok) {
        this.showErrorMsg = true;
      }
    })
  }

  postClaGithubOrganization() {
    let organization = {
      organization_sfid: this.claProjectId,
      organization_name: this.form.value.orgName
    };
    this.claService.postGithubOrganization(organization).subscribe((response) => {
      this.responseErrors = [];

      if (response.errors) {
        this.form.controls['orgName'].setErrors({ incorrect: true });

        for (let errorKey in response.errors) {
          this.responseErrors.push(response.errors[errorKey]);
        }
      } else {
        this.dismiss(true);
      }
    });
  }

  dismiss(data = false) {
    this.viewCtrl.dismiss(data);
  }
}
