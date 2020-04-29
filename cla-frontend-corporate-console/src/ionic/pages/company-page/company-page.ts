// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import {AlertController, IonicPage, ModalController, NavController, NavParams} from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { ClaCompanyModel } from '../../models/cla-company';
import { ClaUserModel } from '../../models/cla-user';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';
import { ColumnMode, SelectionType, SortType } from '@swimlane/ngx-datatable';

@Restricted({
  roles: ['isAuthenticated']
})
@IonicPage({
  segment: 'company/:companyId'
})
@Component({
  selector: 'company-page',
  templateUrl: 'company-page.html'
})
export class CompanyPage {
  pendingRequests: any;
  companyId: string;
  company: ClaCompanyModel;
  manager: ClaUserModel;

  ColumnMode = ColumnMode;
  SelectionType = SelectionType;
  SortType = SortType;

  companySignatures: any[];
  projects: any[];
  loading: any;
  invites: any;

  data: any;
  columns: any[];
  rows: any[] = [];
  allSignatures: any[];
  userEmail: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public alertCtrl: AlertController,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private rolesService: RolesService // for @Restricted
  ) {
    this.companyId = navParams.get('companyId');
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companySignatures: true,
      invites: true,
      projects: true
    };
    this.company = new ClaCompanyModel();
    this.pendingRequests = [];
    this.userEmail = localStorage.getItem('user_email');

    this.data = {};
    this.columns = [
      { prop: 'ProjectName' },
      { prop: 'ProjectManagers' },
      { prop: 'Status' },
      { prop: 'PendingRequets' },
      { prop: 'Approved List' }
    ];
  }

  ngOnInit() {
    this.getCompany();
    this.getCompanySignatures();
    this.getInvites();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe((response) => {
      this.company = response;
      this.getUser(this.company.company_manager_id);
    });
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe((response) => {
      this.manager = response;
    });
  }

  getCompanySignatures() {
    this.loading.companySignatures = true;
    this.loading.projects = true;

    this.claService.getCompanySignatures(this.companyId).subscribe(
      (response) => {
        this.loading.companySignatures = false;

        if (response.resultCount > 0) {
          this.companySignatures = response.signatures.filter((signature) => signature.signatureSigned === true);
          if (this.companySignatures.length <= 0) {
            this.loading.projects = false;
          }
          for (let signature of this.companySignatures) {
            this.getProject(signature);
          }
        } else {
          this.loading.projects = false;
        }
      },
      (exception) => {
        this.loading.companySignatures = false;
        this.loading.projects = false;
      }
    );
  }

  getProject(signature) {
    this.claService.getProject(signature.projectID).subscribe((response) => {
      this.mapProjects(response, signature.signatureACL);
    });
  }

  mapProjects(projectDetail, signatureACL) {
    if (projectDetail) {
      this.claService.getProjectWhitelistRequest(this.companyId, projectDetail.project_id, "pending").subscribe((res) => {
        let pendingContributorRequests = [];
        this.loading.projects = false;
        if (res.list.length > 0) {
          pendingContributorRequests = res.list.filter((r) => {
            return r.projectId === projectDetail.project_id
          })
        }
        this.rows.push({
          ProjectID: projectDetail.project_id,
          ProjectName: projectDetail.project_name !== undefined ? projectDetail.project_name : '',
          ProjectManagers: signatureACL,
          Status: this.getStatus(this.companySignatures),
          PendingContributorRequests: pendingContributorRequests.length,
        });
        this.rows.sort((a, b) => {
          return a.ProjectName.toLowerCase().localeCompare(b.ProjectName.toLowerCase());
        });
      })
    }
  }

  onSelect(projectId) {
    this.openProjectPage(projectId);
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      companyId: this.companyId,
      projectId: projectId
    });
  }

  viewCLAManager(row) {
    let modal = this.modalCtrl.create('ViewCLAManagerModal', {
      'managers': row.ProjectManagers,
      'ProjectName': row.ProjectName
    });
    modal.present();
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'email',
      company: this.company
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'domain',
      company: this.company
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openProjectsCclaSelectModal() {
    let modal = this.modalCtrl.create('ProjectsCclaSelectModal', {
      company: this.company,
      companyId: this.companyId
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  getInvites() {
    this.claService.getCompanyInvites(this.companyId, "pending").subscribe((response) => {
      this.invites = response;
      this.loading.invites = false;
    });
  }

  acceptCompanyInvite(invite) {
    let alert = this.alertCtrl.create({
      subTitle: `Accept Request - Confirmation`,
      message: 'This will dismiss this pending request to join the company and send the company ' +
        'employee an email confirming that they have access.<br/><br/>' +
        'Are you sure?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Accept',
          handler: () => {
            this.claService.approveCompanyInvite(this.companyId, invite.inviteId).subscribe((response) => {
              this.getInvites();
            });
          }
        }
      ]
    });
    alert.present();
  }

  declineCompanyInvite(invite) {
    let alert = this.alertCtrl.create({
      subTitle: `Reject Request - Confirmation`,
      message: 'This will dismiss this pending request to join the company and send the company ' +
        'employee an email indicating that their request was rejected.<br/><br/>' +
        'Are you sure?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Accept',
          handler: () => {
            this.claService.rejectCompanyInvite(this.companyId, invite.inviteId).subscribe((response) => {
              this.getInvites();
            });
          }
        }
      ]
    });
    alert.present();
  }

  getStatus(signatures) {
    for (let i = 0; i < signatures.length; i++) {
      return (this.checkStatusOfSignature(signatures[i].signatureACL, this.userEmail))
    }
  }

  checkStatusOfSignature(signatureACL, userEmail) {
    for (let i = 0; i < signatureACL.length; i++) {
      if (signatureACL[i].lfEmail === userEmail) {
        return 'CLA Manager';
      }
    }

    for (let i = 0; i < this.invites.length; i++) {
      if (this.invites[i].userEmail === userEmail) {
        return 'Pending';
      }
    }
    return 'Request Access'
  }
}
