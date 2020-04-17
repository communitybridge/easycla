// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { ClaService } from '../../../services/cla.service';
import { SFProjectModel } from '../../../models/sfdc-project-model';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser']
})
@IonicPage({
  segment: 'project/:projectId'
})
@Component({
  selector: 'project',
  templateUrl: 'project.html'
})
export class ProjectPage {
  selectedProject: any;
  projectId: string;

  project = new SFProjectModel();

  loading: any;
  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    private claService: ClaService
  ) {
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {

  }

  getSFDCProject(projectId) {
    this.claService.getProjectFromSFDC(projectId).subscribe((response) => {
      if (response) {
        this.project = response;
        this.loading.project = false;
      }
    });
  }

  getDefaults() {
    this.loading = {
      project: true
    };
    this.project = {
      id: 'mock_project_id',
      name: 'Mock Project AOSP',
      logoRef: 'mocklogo.com',
      description: 'description'
    };
    this.sort = {
      alert: {
        arrayProp: 'alert',
        sortType: 'text',
        sort: null
      },
      company: {
        arrayProp: 'org.name',
        sortType: 'text',
        sort: null
      },
      product: {
        arrayProp: 'product',
        sortType: 'text',
        sort: null
      },
      status: {
        arrayProp: 'invoices[0].status',
        sortType: 'text',
        sort: null
      },
      dues: {
        arrayProp: 'annualDues',
        sortType: 'number',
        sort: null
      },
      renewal: {
        arrayProp: 'renewalDate',
        sortType: 'date',
        sort: null
      }
    };
  }
}
