// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { AlertController, IonicPage, ModalController, NavController, NavParams } from 'ionic-angular';
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
  claManagerRequests: any;

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
      claManagerRequests: true,
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
    this.getCompanyInvites();
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
    this.loading.projects = true;
    this.claService.getCompanySignatures(this.companyId).subscribe(
      (response) => {
        if (response.resultCount > 0) {
          this.companySignatures = response.signatures.filter((signature) => signature.signatureSigned === true);
          this.loading.projects = this.companySignatures.length <= 0 ? false : true;
          for (let signature of this.companySignatures) {
            this.getProject(signature);
          }
        } else {
          this.loading.projects = false;
        }
      },
      (exception) => {
        this.loading.projects = false;
      }
    );
  }

  getProject(signature) {
    this.claService.getProject(signature.projectID).subscribe((response) => {
      this.rows.push({
        ProjectID: response.project_id,
        ProjectName: response.project_name !== undefined ? response.project_name : '',
        ProjectManagers: signature.signatureACL,
        Status: this.getStatus(response.project_id, this.rows.length),
        PendingContributorRequests: this.getPendingContributorRequests(response.project_id, this.rows.length),
        PendingCLAManagerRequests: this.getPendingCLAManagerRequests(response.project_id, this.rows.length),
      });
    });
  }

  sortData() {
    this.rows.sort((a, b) => {
      return a.ProjectName.toLowerCase().localeCompare(b.ProjectName.toLowerCase());
    });
  }

  getPendingContributorRequests(projectId, index) {
    this.claService.getProjectWhitelistRequest(this.companyId, projectId, "pending").subscribe((res) => {
      let pendingContributorRequests = [];
      if (res.list.length > 0) {
        pendingContributorRequests = res.list.filter((r) => {
          return r.projectId === projectId
        })
      }
      this.rows[index].PendingContributorRequests = pendingContributorRequests.length;
    });
    return '-';
  }

  getPendingCLAManagerRequests(projectId, index) {
    this.claService.getCLAManagerRequests(this.companyId, projectId).subscribe((response) => {
      let numCLAManagerRequests = 0;
      if (response.requests != null && response.requests.length > 0) {
        numCLAManagerRequests = response.requests.filter((req) => req.status === 'pending').length;
      }
      this.rows[index].PendingCLAManagerRequests = numCLAManagerRequests;
    });
    return '-';
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

  /**
   * Gets the pending company invites for this company.
   */
  getCompanyInvites() {
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
              this.getCompanyInvites();
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
              this.getCompanyInvites();
            });
          }
        }
      ]
    });
    alert.present();
  }

  getStatus(projectId, index) {
    for (let i = 0; i < this.companySignatures.length; i++) {
      const currentProjectId = this.companySignatures[i].projectID;
      if (currentProjectId === projectId) {
        const signatureACL = this.companySignatures[i].signatureACL;
        return (this.checkStatusOfSignature(signatureACL, projectId, index))
      }
    }
  }

  checkStatusOfSignature(signatureACL, projectId, index) {
    const isManager = this.isCLAManger(signatureACL);
    if (!isManager) {
      let status = 'Request Access';
      this.claService.getCLAManagerRequests(this.companyId, projectId).subscribe((response) => {
        if (response.requests != null && response.requests.length > 0) {
          const userId = localStorage.getItem('userid');
          let pendingCLAManagerRequests = response.requests.filter((req) => req.status === 'pending' && req.userID === userId);
          if (pendingCLAManagerRequests.length > 0) {
            status = 'Pending';
          }
        }
        this.rows[index].Status = status;
      }, (error) => {
        this.rows[index].Status = '-';
        console.log(`error loading cla manager requests: ${error}`);
      })
    } else {
      return 'CLA Manager';
    }
    return '-';
  }

  isCLAManger(signatureACL) {
    for (let i = 0; i < signatureACL.length; i++) {
      if (signatureACL[i].lfEmail === this.userEmail) {
        return true;
      }
    }
    return false;
  }


  claManagerRequest(companyID: string, companyName: string, projectID: string, projectName: string, index) {
    let alert = this.alertCtrl.create({
      subTitle: `CLA Manager Request - Confirmation`,
      message: `This will send an email to all the CLA Managers for ${companyName} associated with project ${projectName}.` +
        'The email will instruct the CLA Managers on how to log into the EasyCLA Corporate Console and review this request.' +
        '<br/><br/>' +
        'Are you sure you want to send the request?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Send Request',
          handler: () => {
            const userId = localStorage.getItem('userid');
            const userEmail = localStorage.getItem('user_email');
            const userName = localStorage.getItem('user_name');
            this.claService.createCLAManagerRequest(companyID, projectID, userName, userEmail, userId).subscribe((response) => {
              this.rows[index].Status = 'Pending';
            });
          }
        }
      ]
    });
    alert.present();
  }
}
