// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {AlertController, IonicPage, ModalController, NavController, NavParams} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';
import {ClaCompanyModel} from '../../models/cla-company';
import {ClaUserModel} from '../../models/cla-user';
import {RolesService} from '../../services/roles.service';
import {Restricted} from '../../decorators/restricted';
import {ColumnMode, SelectionType, SortType} from '@swimlane/ngx-datatable';

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
        this.claService.getCLAManagerRequests(this.companyId, projectDetail.project_id).subscribe((response) => {
          let numCLAManagerRequests = 0;
          if (response.requests != null && response.requests.length > 0) {
            numCLAManagerRequests = response.requests.filter((req) => req.status === 'pending').length;
          }
          this.rows.push({
            ProjectID: projectDetail.project_id,
            ProjectName: projectDetail.project_name !== undefined ? projectDetail.project_name : '',
            ProjectManagers: signatureACL,
            Status: this.getStatus(this.companySignatures),
            PendingContributorRequests: (pendingContributorRequests != null) ? pendingContributorRequests.length : 0,
            PendingCLAManagerRequests: numCLAManagerRequests,
          });
          this.rows.sort((a, b) => {
            return a.ProjectName.toLowerCase().localeCompare(b.ProjectName.toLowerCase());
          });
        })
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

  getStatus(signatures) {
    for (let i = 0; i < signatures.length; i++) {
      return (this.checkStatusOfSignature(signatures[i].signatureACL, signatures[i].projectID, this.userEmail))
    }
  }

  checkStatusOfSignature(signatureACL, projectID: string, userEmail: string) {
    for (let i = 0; i < signatureACL.length; i++) {
      if (signatureACL[i].lfEmail === userEmail) {
        return 'CLA Manager';
      }
    }

    let status = 'Request Access';
    this.claService.getCLAManagerRequests(this.companyId, projectID).subscribe((response) => {
      console.log('response:');
      console.log(response);
      if (response.requests != null && response.requests.length > 0) {
        const userId = localStorage.getItem('userid');
        console.log(`Testing ${userId}`);
        let pendingCLAManagerRequests = response.requests.filter((req) => req.status === 'pending' && req.userID === userId);
        console.log(`pending requests:`);
        console.log(pendingCLAManagerRequests);
        if (pendingCLAManagerRequests.length > 0) {
          console.log('setting pending status');
          status = 'Pending';
        }
      }
    }, (error) => {
      console.log(`error loading cla manager requests: ${error}`);
    })

    console.log(`returnging status: ${status}`);
    return status;
  }


  claManagerRequest(companyID: string, companyName: string, projectID: string, projectName: string) {
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
            });
          }
        }
      ]
    });
    alert.present();
  }
}
