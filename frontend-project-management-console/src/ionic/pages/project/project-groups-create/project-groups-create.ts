// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { DomSanitizer} from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project/:projectId/groups/create'
// })
@Component({
  selector: 'project-groups-create',
  templateUrl: 'project-groups-create.html'
})

export class ProjectGroupsCreatePage {

  projectId: string;
  keysGetter;
  projectPrivacy;

  groupName: string;
  groupDescription: string;
  groupPrivacy = [];

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  group: any;
  projectGroups: any;
  approveMembers: boolean;
  restrictPosts: boolean;
  approvePosts: boolean;
  allowUnsubscribed: boolean;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
    private formBuilder: FormBuilder,
  ) {
    this.projectId = navParams.get('projectId');

    this.form = formBuilder.group({
      groupName:[this.groupName, Validators.compose([Validators.minLength(3), Validators.pattern(/^\S*$/), Validators.required])],
      groupDescription:[this.groupDescription, Validators.compose([Validators.minLength(9), Validators.required])],
      groupPrivacy:[this.groupPrivacy, Validators.compose([Validators.required])],
      approveMembers: [this.approveMembers],
      restrictPosts: [this.restrictPosts],
      approvePosts: [this.approvePosts],
      allowUnsubscribed: [this.allowUnsubscribed]
    });

  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    this.getProjectGroups();
    this.getGroupPrivacy();
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        console.log(response);
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
      console.log(response);
    });
  }

  getGroupPrivacy() {
    this.groupPrivacy = [];
    // TODO Implement CINCO side
    // this.cincoService.getGroupPrivacy(this.projectId).subscribe(response => {
    //   this.groupPrivacy = response;
    // });
    this.groupPrivacy = [
      {
        value: "sub_group_privacy_none",
        description: "Group listed and archive publicly viewable"
      },
      {
        value: "sub_group_privacy_archives",
        description: "Group listed and archive privately viewable by members"
      },
      {
        value: "sub_group_privacy_unlisted",
        description: "Group hidden and archive privately viewable by members"
      }
    ];
  }

  submitGroup() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    // CINCO / Groups.io Flags
    // allow_unsubscribed - Allow unsubscribed members to post
    // approve_members    - Members required approval to join
    // approve_posts      - Posts require approval
    // restrict_posts     - Posts are restricted to moderators only
    this.group = {
      projectId: this.projectId,
      name: this.form.value.groupName,
      description: this.form.value.groupDescription,
      privacy: this.groupPrivacy[this.form.value.groupPrivacy].value,
      allow_unsubscribed: this.form.value.allowUnsubscribed,
      approve_members: this.form.value.approveMembers,
      approve_posts: this.form.value.approvePosts,
      restrict_posts: this.form.value.restrictPosts
    };
    console.log(this.group);
    this.cincoService.createProjectGroup(this.projectId, this.group).subscribe(response => {
      this.currentlySubmitting = false;
      this.navCtrl.setRoot('ProjectGroupsPage', {
        projectId: this.projectId
      });
    });
  }

  getGroupsList() {
    this.navCtrl.setRoot('ProjectGroupsPage', {
      projectId: this.projectId
    });
  }

}
