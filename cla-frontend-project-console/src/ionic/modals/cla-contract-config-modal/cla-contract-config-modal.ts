// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { Events, IonicPage, NavController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { PlatformLocation } from '@angular/common';
import { generalConstants } from '../../constants/general';

@IonicPage({
  segment: 'cla-contract-config-modal'
})
@Component({
  selector: 'cla-contract-config-modal',
  templateUrl: 'cla-contract-config-modal.html'
})
export class ClaContractConfigModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  projectId: string;
  claProject: any;
  newClaProject: boolean;
  loading: boolean;
  errorMessage: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    private claService: ClaService,
    public events: Events,
    private location: PlatformLocation
  ) {
    this.projectId = this.navParams.get('projectId');
    this.claProject = this.navParams.get('claProject');
    this.getDefaults();
    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
    this.form = formBuilder.group({
      name: [this.claProject.projectName, Validators.compose(
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(100),
          Validators.pattern(new RegExp(generalConstants.ALLOW_ALL_LANGUAGES))
        ])],
      description: [this.claProject.projectDescription, Validators.compose(
        [
          Validators.maxLength(255),
          Validators.pattern(new RegExp(generalConstants.ALLOW_ALL_LANGUAGES))
        ])],
      ccla: [this.claProject.projectCCLAEnabled],
      cclaAndIcla: [this.claProject.projectCCLARequiresICLA],
      icla: [this.claProject.projectICLAEnabled]
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  getDefaults() {
    this.newClaProject = false; // we assume we have an existing cla project
    // if claProject is not passed
    if (!this.claProject) {
      this.newClaProject = true; // change to creating new project
      this.claProject = {
        projectExternalID: this.projectId,
        projectName: '',
        projectDescription: '',
        projectCCLAEnabled: false,
        projectCCLARequiresICLA: false,
        projectICLAEnabled: false
      };
    }
  }

  ngOnInit() {
  }

  submit() {
    this.submitAttempt = true;
    this.errorMessage = '';
    this.currentlySubmitting = true;
    if (this.isFormValid()) {
      if (this.newClaProject) {
        this.postProject();
      } else {
        this.putProject();
      }
    } else {
      this.currentlySubmitting = false;
    }
  }

  isFormValid() {
    return this.form.valid && (this.form.controls.ccla.value || this.form.controls.icla.value)
  }

  checkMandatory() {
    if (!this.form.controls.ccla.value || !this.form.controls.icla.value) {
      this.form.controls['cclaAndIcla'].setValue(false);
    }
  }

  postProject() {
    this.loading = true;
    let claProject = {
      projectExternalID: this.claProject.projectExternalID,
      projectName: this.form.value.name,
      projectDescription: this.form.value.description,
      projectCCLAEnabled: this.form.value.ccla,
      projectACL: [localStorage.getItem('userid')],
      projectCCLARequiresICLA: this.form.value.cclaAndIcla,
      projectICLAEnabled: this.form.value.icla
    };
    this.claService.postProject(claProject).subscribe((response) => {
      this.loading = false;
      this.dismiss(true);
    }, (error) => {
      if (error) {
        this.loading = false;
        this.errorMessage = JSON.parse(error._body).message;
      }
    });
  }

  putProject() {
    // rebuild the claProject object from existing data and form data
    this.loading = true;
    let claProject = {
      projectID: this.claProject.projectID,
      projectExternalID: this.claProject.projectExternalID,
      projectName: this.form.value.name,
      projectDescription: this.form.value.description,
      projectCCLAEnabled: this.form.value.ccla,
      projectCCLARequiresICLA: this.form.value.cclaAndIcla,
      projectICLAEnabled: this.form.value.icla
    };
    this.claService.putProject(claProject).subscribe(
      (res) => {
        this.loading = false;
        this.dismiss(true);
      },
      (error) => {
        this.loading = false;
        if (!error.ok) {
          this.errorMessage = JSON.parse(error._body).Message;
        }
      });
  }

  dismiss(data?) {
    this.viewCtrl.dismiss(data);
  }

  clearError(event) {
    this.errorMessage = '';
  }
}
