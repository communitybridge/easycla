// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAdmin']
})
// @IonicPage({
//   segment: 'activity-log'
// })
@Component({
  selector: 'activity-log',
  templateUrl: 'activity-log.html'
})
export class ActivityLogPage {
  allProjects: any;
  selectedProject: any;
  events: any;
  loading: any;
  users: any;
  expand: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private rolesService: RolesService // for @Restricted
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.events = [];
    this.loading = {
      events: true
    };
    this.expand = {};
    this.users = {};
  }

  ngOnInit() {
    this.getProjects();
  }

  getProjects() {
    this.cincoService.getAllProjects().subscribe(response => {
      this.allProjects = response;
      this.selectedProject = this.allProjects[0];
      this.projectSelectChanged();
    });
  }

  projectSelectChanged() {
    let projectId = this.selectedProject.id;
    this.events = [];
    this.loading.events = true;
    this.getEvents(projectId);
  }

  getEvents(projectId) {
    this.cincoService.getEventsForProject(projectId).subscribe(response => {
      if (response) {
        this.events = response;
        for (let event of this.events) {
          this.getUser(event.userId);
        }
        this.loading.events = false;
      }
    });
  }

  getUser(userId) {
    if (!this.users[userId]) {
      this.users[userId] = true; //placeholder
      this.cincoService.getUser(userId).subscribe(response => {
        if (response) {
          this.users[userId] = response;
        }
      });
    }
  }

  toggle(index) {
    this.expand[index] = !this.expand[index];
  }
}
