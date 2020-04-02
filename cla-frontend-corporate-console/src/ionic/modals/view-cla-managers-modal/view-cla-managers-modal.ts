// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { ViewController, IonicPage, NavParams } from 'ionic-angular';
import { FormBuilder } from '@angular/forms';
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
  filteredManagers: any;
  ProjectName: string;

  constructor(
    public viewCtrl: ViewController,
    public navParams: NavParams,
    public formBuilder: FormBuilder,
    location: PlatformLocation,
  ) {
    this.getDefaults();
    location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
  }

  getDefaults() {
    this.managers = this.navParams.get('managers');
    this.filteredManagers = this.managers;
    this.ProjectName = this.navParams.get('ProjectName');
  }

  trimCharacter(text, length) {
    return text.length > length ? text.substring(0, length) + '...' : text;
  }

  dismiss(data = false) {
    this.viewCtrl.dismiss(data);
  }

  searchCLAManager(event) {
    const keyword = event.value.toLowerCase();
    this.filteredManagers = this.managers.filter((manager) => {
      return (manager.username.toLowerCase().indexOf(keyword) >= 0);
    });
  }
}
