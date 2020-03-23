// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input } from '@angular/core';
import { ViewController, IonicPage, NavParams } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { AuthService } from '../../services/auth.service';
import { ClaCompanyModel } from '../../models/cla-company';
import { PlatformLocation } from '@angular/common';

@IonicPage({
  segment: 'view-cla-managers-modal'
})
@Component({
  selector: 'view-cla-managers-modal',
  templateUrl: 'view-cla-managers-modal.html'
})
export class ViewCLAManagerModal {
  managers: any;
  ProjectName: string;

  constructor(
    public viewCtrl: ViewController,
    public navParams: NavParams,
    public formBuilder: FormBuilder,
    private claService: ClaService,
    location: PlatformLocation,
  ) {
    this.getDefaults();
    location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
  }

  getDefaults() {
    this.managers = this.navParams.get('managers');
    this.ProjectName = this.navParams.get('ProjectName');
  }

  dismiss(data = false) {
    this.viewCtrl.dismiss(data);
  }
}
