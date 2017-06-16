import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service';


@IonicPage({
  segment: 'project/:projectId/member/:memberId'
})
@Component({
  selector: 'page-member',
  templateUrl: 'member.html',
  providers: [CincoService],
})
export class MemberPage {
  projectId: any;
  memberId: any;
  member: any;
  memberContacts: any;
  memberContactRoles: any;
  orgProjectMemberships: any;
  orgProjectMembershipsFiltered: any;

  constructor(
    public navCtrl: NavController,
    public modalCtrl: ModalController,
    public navParams: NavParams,
    private cincoService: CincoService
  ) {
    // If we navigated to this page, we will have an item available as a nav param
    this.projectId = navParams.get('projectId');
    this.memberId = navParams.get('memberId');

    this.getDefaults();
  }

  ngOnInit() {
    // use selectedProject and selectedMember to get data from CINCO and populate this.contacts
    this.getMemberContactRoles();
    this.getMember(this.projectId, this.memberId);
    this.getMemberContacts(this.projectId, this.memberId);
  }

  getMember(projectId, memberId) {
    this.cincoService.getMember(projectId, memberId).subscribe(response => {
      if(response) {
        this.member = response;
        this.getOrganizationProjectMemberships(this.member.org.id);
      }
    });
  }

  getMemberContacts(projectId, memberId) {
    this.cincoService.getMemberContacts(projectId, memberId).subscribe(response => {
      if(response) {
        this.memberContacts = response;
      }
    });
  }

  addMemberContact() {
    let modal = this.modalCtrl.create('SearchAddContact', {
      projectId: this.projectId,
      memberId: this.memberId,
      org: this.member.org,
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getMemberContacts(this.projectId, this.memberId);
    });
    modal.present();
  }

  contactSelected(event, contact) {
    let modal = this.modalCtrl.create('ContactUpdate', {
      projectId: this.projectId,
      memberId: this.member.id,
      org: this.member.org,
      contact: contact,
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getMemberContacts(this.projectId, this.memberId);
    });
    modal.present();
  }

  getDefaults() {
    this.member = {
      id: this.memberId,
      contacts: [
        {
          bio: "",
          email: "",
          familyName: "",
          givenName: "",
          headshotRef: "",
          id: "",
          phone: "",
          type: "",
        }
      ],
      invoices: [],
      org: {
        addresses: [],
        id: "",
        logoRef: "",
        name: "",
        phone: "",
        url: "",
      },
      projectId: "",
      renewalDate: "",
      startDate: "",
      tier: {
        qualifier: "",
        type: "",
      }
    };
    this.memberContactRoles = {};
  }

  getMemberContactRoles() {
    this.cincoService.getMemberContactRoles().subscribe(response => {
      if(response) {
        this.memberContactRoles = response;
      }
    });
  }

  getOrganizationProjectMemberships(organizationId) {
    this.cincoService.getOrganizationProjectMemberships(organizationId).subscribe(response => {
      if(response) {
        this.orgProjectMemberships = response;
        this.orgProjectMembershipsFiltered = this.orgProjectMemberships.filter((item, index) => index < 4 );
      }
    });
  }

  openMembershipsModal() {
    let modal = this.modalCtrl.create('MembershipsModal', {
      orgName: this.member.org.name,
      memberships: this.orgProjectMemberships,
    });
    modal.present();
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      projectId: projectId,
    });
  }

}
