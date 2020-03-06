// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {Events, IonicPage, NavController, NavParams, ViewController} from 'ionic-angular';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {ClaService} from '../../services/cla.service';
import {CLAGroupModel} from '../../models/cla-group-model';

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

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    private claService: ClaService,
    public events: Events
  ) {
    this.projectId = this.navParams.get('projectId');
    this.claProject = this.navParams.get('claProject');
    this.getDefaults();
    this.form = formBuilder.group({
      name: [this.claProject.project_name, Validators.compose([Validators.required])],
      ccla: [this.claProject.project_ccla_enabled],
      cclaAndIcla: [this.claProject.project_ccla_requires_icla_signature],
      icla: [this.claProject.project_icla_enabled]
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
        project_external_id: this.projectId,
        project_name: '',
        project_ccla_enabled: false,
        project_ccla_requires_icla_signature: false,
        project_icla_enabled: false
      };
    }
  }

  ngOnInit() {
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    if (this.newClaProject) {
      this.createCLAGroup();
    } else {
      this.updateCLAGroup();
    }
  }

  checkMandatory(value: boolean = true) {
    this.form.controls['cclaAndIcla'].setValue(value);
  }

  createCLAGroup() {
    let claGroupModel = new CLAGroupModel(
      null,
      this.form.value.name,
      this.claProject.project_external_id,
      [localStorage.getItem('userid')], // Current user
      this.form.value.icla,
      this.form.value.ccla,
      this.form.value.cclaAndIcla
    );
    this.claService.createCLAGroup(claGroupModel).subscribe((response) => {
      this.dismiss();
    });
  }

  updateCLAGroup() {
    // rebuild the claProject object from existing data and form data
    let claGroupModel = new CLAGroupModel(
      this.claProject.project_id,
      this.form.value.name,
      this.claProject.project_external_id,
      null,// [localStorage.getItem('userid')], // Current user
      this.form.value.icla,
      this.form.value.ccla,
      this.form.value.cclaAndIcla
    );
    this.claService.updateCLAGroup(claGroupModel).subscribe((response) => {
      this.dismiss();
    });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
