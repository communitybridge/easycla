// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { SortService } from '../../services/sort.service';

@IonicPage({
  segment: 'memberships-modal'
})
@Component({
  selector: 'memberships-modal',
  templateUrl: 'memberships-modal.html',
  providers: [SortService]
})
export class MembershipsModal {
  orgName: string;
  memberships: any;
  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private sortService: SortService
  ) {
    this.orgName = navParams.get('orgName');
    this.memberships = navParams.get('memberships');
    this.getDefaults();
  }

  ngOnInit() {}

  getDefaults() {
    this.sort = {
      project: {
        arrayProp: 'projectName',
        sortType: 'text',
        sort: null
      },
      class: {
        arrayProp: 'product',
        sortType: 'text',
        sort: null
      },
      status: {
        arrayProp: 'invoices[0].status',
        sortType: 'text',
        sort: null
      },
      renewal: {
        arrayProp: 'renewalDate',
        sortType: 'date',
        sort: null
      }
    };
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      projectId: projectId
    });
  }

  sortProjects(prop) {
    this.sortService.toggleSort(this.sort, prop, this.memberships);
  }
}
