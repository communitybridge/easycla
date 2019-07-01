// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { DomSanitizer} from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';
import { AlertController } from 'ionic-angular';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project/:projectId/groups'
// })
@Component({
  selector: 'project-groups',
  templateUrl: 'project-groups.html'
})

export class ProjectGroupsPage {

  projectId: string;

  groupName: string;
  groupDescription: string;
  groupPrivacy = [];
  groupRequiresApproval = [];

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  group: any;
  projectGroups: any;
  allGroupsWithParticipants: any[] = [];
  groupParticipants: any;

  participantName: any;
  participantEmail: any;

  expand: any;
  isTrue:boolean = null;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
    private formBuilder: FormBuilder,
    private alertCtrl: AlertController
  ) {
    this.projectId = navParams.get('projectId');

    this.form = formBuilder.group({
      groupName:[this.groupName, Validators.compose([Validators.required])],
      groupDescription:[this.groupDescription],
      groupPrivacy:[this.groupPrivacy],
      groupRequiresApproval:[this.groupRequiresApproval],
    });

    this.form = formBuilder.group({
      participantName:[this.participantName],
      participantEmail:[this.participantEmail, Validators.compose([Validators.required])],
    });

  }

  async ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.expand = {};
    this.projectGroups = [];
    this.allGroupsWithParticipants = [];
    this.form.reset();
    this.getProjectGroups();
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        if (!response.mailingGroup) {
          console.log("no mailingGroup");
          console.log("creating a new mailingGroup");
          this.cincoService.createMainProjectGroup(this.projectId).subscribe(response => {
            console.log("new mailingGroup");
            console.log(response);
          });
        }
      }
    });
  }

  getProjectGroups() {
    this.cincoService.getAllProjectGroups(this.projectId).subscribe(response => {
      this.projectGroups = response;
      for(let eachProject of this.projectGroups) {
        this.getAllGroupParticipants(eachProject.name);
      }
    });
  }

  getAllGroupParticipants(groupName){
    let group = {
      info: [],
      participants: []
    };
    this.cincoService.getProjectGroup(this.projectId, groupName).subscribe(response => {
      group.info = response;
      this.cincoService.getAllGroupParticipants(this.projectId, groupName).subscribe(response => {
        group.participants = response;
        this.allGroupsWithParticipants.push(group);
      });
    });
  }

  submitGroup() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.group = {
      projectId: this.projectId,
      name: this.form.value.groupName,
      description: this.form.value.groupDescription,
      privacy: this.groupPrivacy[this.form.value.groupPrivacy].value,
      requires_approval: this.groupRequiresApproval[this.form.value.groupRequiresApproval].value,
    };
    console.log(this.group);
    this.cincoService.createProjectGroup(this.projectId, this.group).subscribe(response => {
      this.currentlySubmitting = false;
      this.navCtrl.setRoot('ProjectGroupsPage', {
        projectId: this.projectId
      });
    });
  }

  goCreateGroupsPage() {
    this.navCtrl.setRoot('ProjectGroupsCreatePage', {
      projectId: this.projectId
    });
  }

  toggle(index) {
    this.expand[index] = !this.expand[index];
  }

  submitParticipant(groupName) {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let participant = [{
        address: this.form.value.participantEmail,
        name: this.form.value.participantName
    }];
    this.cincoService.addGroupParticipant(this.projectId, groupName, participant).subscribe(response => {
      this.currentlySubmitting = false;
      this.getDefaults();
    });
  }

  removeGroupParticipant(groupName, participantEmail){
    let participant = [{
        address: participantEmail
    }];
    this.presentRemoveConfirm((confirm) => {
      if(confirm) {
        this.cincoService.removeGroupParticipant(this.projectId, groupName, participant).subscribe(response => {
          // CINCO and Groups.io sync takes a while
          let refreshData = setTimeout( () => {
            this.getDefaults()
          }, 2000);
        });
      }

    });
  }

  removeProjectGroup(groupName){
    this.presentRemoveConfirm((confirm) => {
      if(confirm) {
        this.cincoService.removeProjectGroup(this.projectId, groupName).subscribe(response => {
          console.log(response);
          // CINCO and Groups.io sync takes a while
          let refreshData = setTimeout( () => {
            this.getDefaults()
          }, 2000);
        });
      }
    });
  }

  presentRemoveConfirm(callback) {
    let alert = this.alertCtrl.create({
      title: 'Please Confirm',
      message: 'Are you sure you want to remove this?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            callback(false);
          }
        },
        {
          text: 'Remove',
          handler: () => {
            callback(true);
          }
        }
      ]
    });
    alert.present();
  }

}
