// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';
import { NavController, ModalController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { RolesService } from '../../services/roles.service';

@Component({
  selector: 'member-header',
  templateUrl: 'member-header.html'
})
export class MemberHeaderComponent {
  @Input('memberId')
  private memberId: string;

  @Input('projectId')
  private projectId: string;

  can: any;

  loading: any;

  member: any;

  memberships: any;

  membershipsFiltered: any;

  userRoles: any;

  constructor(
    private navCtrl: NavController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    public rolesService: RolesService
  ) {
    this.getDefaults();
  }

  ngOnInit() {
    this.getMember();

    this.rolesService.getUserRolesPromise().then(userRoles => {
      this.userRoles = userRoles;
      this.can.viewPartnerships = !userRoles.isStaffInc;
      this.getMember();
    });
  }

  getDefaults() {
    this.can = {
      viewPartnerships: false
    };
    this.loading = {
      member: true,
      memberships: true
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

    this.memberships = [];
  }

  getMember() {
    this.cincoService.getMember(this.projectId, this.memberId).subscribe(response => {
      if (response) {
        this.member = response;
        this.loading.member = false;
        if (!this.userRoles.isStaffInc) {
          this.getOrganizationProjectMemberships(this.member.org.id);
        }
      }
    });
  }

  openMembershipsModal() {
    let modal = this.modalCtrl.create('MembershipsModal', {
      orgName: this.member.org.name,
      memberships: this.memberships
    });
    modal.present();
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      projectId: projectId
    });
  }

  getOrganizationProjectMemberships(organizationId) {
    this.cincoService.getOrganizationProjectMemberships(organizationId).subscribe(response => {
      if (response) {
        this.memberships = response;
        this.loading.memberships = false;
        this.membershipsFiltered = this.memberships.filter((item, index) => index < 2);
      }
    });
  }
}
