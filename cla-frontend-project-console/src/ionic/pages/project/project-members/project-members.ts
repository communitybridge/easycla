// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';
import { ProjectModel } from '../../../models/project-model';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project/:projectId/members'
// })
@Component({
  selector: 'project-members',
  templateUrl: 'project-members.html',
})
export class ProjectMembersPage {
  selectedProject: any;
  projectId: string;

  project = new ProjectModel();

  loading: any;
  sort: any;

  contracts: any;

  iclaUploadInfo: any;
  cclaUploadInfo: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.selectedProject = navParams.get('project');
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.project = response;
        // This is to refresh an image that have same URL
        if(this.project.config.logoRef) { this.project.config.logoRef += "?" + new Date().getTime(); }
        this.loading.project = false;
      }
    });
  }

  memberSelected(event, memberId) {
    this.navCtrl.push('MemberPage', {
      projectId: this.projectId,
      memberId: memberId,
    });
  }

  getDefaults() {
    this.loading = {
      project: true,
    };
    this.project = {
      id: "",
      name: "Project",
      description: "Description",
      managers: "",
      members: [],
      status: "",
      category: "",
      sector: "",
      url: "",
      logoRef: "",
      startDate: "",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: {
        address: {
          administrativeArea: "",
          country: "",
          localityName: "",
          postalCode: "",
          thoroughfare: ""
        },
        type: ""
      },
      config: {
        logoRef: ""
      }
    };
    this.sort = {
      alert: {
        arrayProp: 'alert',
        sortType: 'text',
        sort: null,
      },
      company: {
        arrayProp: 'org.name',
        sortType: 'text',
        sort: null,
      },
      product: {
        arrayProp: 'product',
        sortType: 'text',
        sort: null,
      },
      status: {
        arrayProp: 'invoices[0].status',
        sortType: 'text',
        sort: null,
      },
      dues: {
        arrayProp: 'annualDues',
        sortType: 'number',
        sort: null,
      },
      renewal: {
        arrayProp: 'renewalDate',
        sortType: 'date',
        sort: null,
      },
    };
  }

  sortMembers(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.project.members,
    );
  }

}
