// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavParams, ViewController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'project-user-management-modal'
})
@Component({
  selector: 'project-user-management-modal',
  templateUrl: 'project-user-management-modal.html'
})
export class ProjectUserManagementModal {
  projectId: string;
  projectName: string;
  userIds: any;
  users: any;
  selectedUsers: any;
  userResults: string[];
  userTerm: any;

  constructor(public navParams: NavParams, public viewCtrl: ViewController, private cincoService: CincoService) {
    this.projectId = this.navParams.get('projectId');
    this.projectName = this.navParams.get('projectName');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProjectConfig();
  }

  getDefaults() {
    this.users = [];
    this.selectedUsers = [];
    this.userResults = [];
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  getProjectConfig() {
    this.cincoService.getProjectConfig(this.projectId).subscribe(response => {
      if (response) {
        this.userIds = response.programManagers;
        this.updateUsers();
      }
    });
  }

  updateUsers() {
    let users = this.userIds;
    this.users = [];
    for (let i = 0; i < users.length; i++) {
      this.appendUser(users[i]);
    }
  }

  appendUser(userId) {
    this.cincoService.getUser(userId).subscribe(response => {
      if (response) {
        this.users.push(response);
      }
    });
  }

  removeUser(userId) {
    let index = this.userIds.indexOf(userId);
    if (index !== -1) {
      this.userIds.splice(index, 1);
      let updatedManagers = JSON.stringify(this.userIds);
      this.cincoService.updateProjectManagers(this.projectId, updatedManagers).subscribe(response => {
        if (response) {
          this.getProjectConfig();
        }
      });
    }
  }

  searchUsers(ev: any) {
    if (this.userTerm) {
      this.cincoService.searchUserTerm(this.userTerm).subscribe(response => {
        if (response) {
          this.userResults = response;
        }
      });
    } else {
      this.userResults = [];
    }
  }

  assignUserPM(id) {
    this.userIds.push(id);
    this.cincoService.updateProjectManagers(this.projectId, this.userIds).subscribe(response => {
      if (response) {
        this.getProjectConfig();
      }
    });
  }
}
