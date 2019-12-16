// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UrlValidator } from '../../validators/url';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'analytics-config-modal'
})
@Component({
  selector: 'analytics-config-modal',
  templateUrl: 'analytics-config-modal.html'
})
export class AnalyticsConfigModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  analyticsUrl: any;
  projectId: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    private cincoService: CincoService
  ) {
    this.getDefaults();
    this.projectId = this.navParams.get('projectId');
    this.form = formBuilder.group({
      analyticsUrl: [this.analyticsUrl, Validators.compose([UrlValidator.isValid])]
    });
  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
  }

  getDefaults() {}

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      return;
    } else {
      this.analyticsUrl = this.form.value.analyticsUrl;
      this.addAnalyticsUrl();
    }
  }

  dismiss() {
    this.viewCtrl.dismiss(this.analyticsUrl);
  }

  addAnalyticsUrl() {
    this.cincoService.getProjectConfig(this.projectId).subscribe((response) => {
      if (response) {
        let updatedConfig = response;
        updatedConfig.analyticsUrl = this.analyticsUrl;
        this.cincoService.editProjectConfig(this.projectId, updatedConfig).subscribe((response) => {
          if (response) {
            this.dismiss();
          }
        });
      }
    });
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe((response) => {
      if (response) {
        let projectConfig = response;
        if (projectConfig.analyticsUrl) {
          this.analyticsUrl = projectConfig.analyticsUrl;
          this.form.patchValue({
            analyticsUrl: this.analyticsUrl
          });
        }
      }
    });
  }
}
