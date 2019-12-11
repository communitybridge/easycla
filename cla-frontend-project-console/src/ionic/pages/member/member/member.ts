// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';
import { MemberModel } from '../../../models/member-model';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser']
})
// @IonicPage({
//   segment: 'project/:projectId/member/:memberId'
// })
@Component({
  selector: 'member',
  templateUrl: 'member.html'
})
export class MemberPage {
  projectId: any;
  memberId: any;
  memberContacts: any;
  memberContactRoles: any;
  orgProjectMemberships: any;
  orgProjectMembershipsFiltered: any;
  loading: any;
  sort: any;

  member = new MemberModel();

  userRoles: any;

  constructor(
    public navCtrl: NavController,
    public modalCtrl: ModalController,
    public navParams: NavParams,
    private sortService: SortService,
    private cincoService: CincoService,
    private rolesService: RolesService,
    private keycloak: KeycloakService
  ) {
    // If we navigated to this page, we will have an item available as a nav param
    this.projectId = navParams.get('projectId');
    this.memberId = navParams.get('memberId');
    this.getDefaults();
  }

  getDefaults() {
    this.userRoles = this.rolesService.userRoles;
    this.loading = {
      member: true,
      projects: true,
      contacts: true
    };
    this.member = {
      id: this.memberId,
      projectId: '',
      projectName: '',
      org: {
        id: '',
        name: '',
        parent: '',
        logoRef: '',
        url: '',
        addresses: ''
      },
      product: '',
      tier: '',
      annualDues: '',
      startDate: '',
      renewalDate: '',
      invoices: ['']
    };

    this.memberContactRoles = {};

    this.sort = {
      role: {
        arrayProp: 'type',
        prop: 'role',
        sortType: 'text',
        sort: null
      },
      name: {
        arrayProp: 'contact.givenName',
        prop: 'name',
        sortType: 'text',
        sort: null
      },
      email: {
        arrayProp: 'contact.email',
        prop: 'email',
        sortType: 'text',
        sort: null
      }
    };
  }

  ngOnInit() {
    this.rolesService.getUserRolesPromise().then(userRoles => {
      this.userRoles = userRoles;

      // use selectedProject and selectedMember to get data from CINCO and populate this.contacts
      this.getMemberContactRoles();
      this.getMember(this.projectId, this.memberId);
      this.getMemberContacts(this.projectId, this.memberId);
    });
  }

  getMemberContactRoles() {
    this.cincoService.getMemberContactRoles().subscribe(response => {
      if (response) {
        this.memberContactRoles = response;
      }
    });
  }

  getMember(projectId, memberId) {
    this.cincoService.getMember(projectId, memberId).subscribe(response => {
      if (response) {
        this.member = response;
        this.loading.member = false;
      }
    });
  }

  getMemberContacts(projectId, memberId) {
    this.cincoService.getMemberContacts(projectId, memberId).subscribe(response => {
      if (response) {
        this.memberContacts = response;
      }
      this.loading.contacts = false;
    });
  }

  addMemberContact() {
    let modal = this.modalCtrl.create('SearchAddContactModal', {
      projectId: this.projectId,
      memberId: this.memberId,
      org: this.member.org
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getMemberContacts(this.projectId, this.memberId);
    });
    modal.present();
  }

  contactSelected(event, contact) {
    let modal = this.modalCtrl.create('ContactUpdateModal', {
      projectId: this.projectId,
      memberId: this.member.id,
      org: this.member.org,
      contact: contact
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getMemberContacts(this.projectId, this.memberId);
    });
    modal.present();
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      projectId: projectId
    });
  }

  sortContacts(prop) {
    this.sortService.toggleSort(this.sort, prop, this.memberContacts);
  }
}
