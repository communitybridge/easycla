// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';

import { NavController, IonicPage, ModalController } from 'ionic-angular';

import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'all-members'
// })
@Component({
  selector: 'all-members',
  templateUrl: 'all-members.html'
})
export class AllMembersPage {
  allProjects: any;
  selectedProject: any;
  projectSectors: any;
  membersList: any;
  membersSelected: any;
  keysGetter: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.getDefaults();
    this.keysGetter = Object.keys;
  }

  getDefaults() {
    this.membersList = {};
    this.membersSelected = {};
    this.projectSectors = [];
  }

  ngOnInit(){
    this.getProjectMembers();
    this.getProjectSectors();
  }

  getProjectSectors() {
    this.cincoService.getProjectSectors().subscribe(response => {
      if(response) {
        this.projectSectors = response;
      }
    });
  }

  getProjectMembers(){
    this.cincoService.getAllProjects().subscribe(response => {
      this.allProjects = response;
      this.selectedProject = this.allProjects[0];
      this.projectSelectChanged();
    });
  }

  projectSelectChanged() {
    let value = this.selectedProject;
    if (value=='all') {
      this.selectProjectMembers(this.allProjects);
    }
    else {
      this.selectProjectMembers([value]);
    }
  }

  selectProjectMembers(projectsArray) {
    this.membersSelected = [];
    for (let i = 0; i < projectsArray.length; i++) {
      let project = projectsArray[i];
      // We already have data for that project's members
      if(this.membersList.hasOwnProperty(project.id)) {
        // Create a reference to it
        this.membersSelected[project.id] = this.membersList[project.id];
      }
      // We need data for this project's members
      else {
        // Look it up and attach it
        this.attachProjectMember(project);
      }
    }
  }

  attachProjectMember(project) {
    this.cincoService.getProjectMembers(project.id).subscribe(response => {
      if(response) {
        let members = response;
        if(members.length == 0){ return; }
        let buildArray = [];
        for (let i = 0; i < members.length; i++) {
          let member = members[i];
          buildArray.push({
            id: member.id,
            logo: member.org.logoRef,
            name: member.org.name,
            projectName: project.name,
            projectId: member.projectId,
            level: member.product,
            invoiceStatus: member.invoices[0].status,
            renewalDate: member.renewalDate,
            projectCategory: project.sector,
          });
        }
        this.membersList[project.id] = buildArray;
        this.membersSelected[project.id] = this.membersList[project.id];
      }
    });
  }

  memberSelected(member) {
    // let modal = this.modalCtrl.create('MemberModal', {
    //
    // });
    // modal.present();
    this.navCtrl.push('MemberPage', {
      projectId: member.projectId,
      memberId: member.id,
    });
  }

}
