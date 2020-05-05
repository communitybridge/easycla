// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {AlertController, IonicPage, ModalController, NavController, NavParams, ViewController} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';
import {ClaCompanyModel} from '../../models/cla-company';
import {ClaUserModel} from '../../models/cla-user';
import {ClaSignatureModel} from '../../models/cla-signature';
import {ClaManager} from '../../models/cla-manager';
import {SortService} from '../../services/sort.service';
import {Restricted} from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated']
})
@IonicPage({
  segment: 'company/:companyId/project/:projectId/:modal'
})
@Component({
  selector: 'project-page',
  templateUrl: 'project-page.html'
})
export class ProjectPage {
  cclaSignature: any;
  employeeSignatures: any[];
  loading: any;
  companyId: string;
  projectId: string;
  projectName: string;
  managers: ClaManager[];
  company: ClaCompanyModel;
  manager: ClaUserModel;
  showModal: any;

  noPendingContributorRequests: boolean;
  allContributorRequests: any[];
  approvedContributorRequests: any[];
  pendingContributorRequests: any[];
  rejectedContributorRequests: any[];

  noPendingCLAManagerRequests: boolean;
  allCLAManagerRequests: any[];
  pendingCLAManagerRequests: any[];
  approvedCLAManagerRequests: any[];
  deniedCLAManagerRequests: any[];

  project: any;
  userEmail: any;
  sort: any;
  canEdit: boolean = false;

  constructor(
    public navParams: NavParams,
    private claService: ClaService,
    public alertCtrl: AlertController,
    public modalCtrl: ModalController,
    private sortService: SortService,
    public viewCtrl: ViewController,
    public navCtrl: NavController
  ) {
    this.userEmail = localStorage.getItem('user_email');
    this.companyId = navParams.get('companyId');
    this.projectId = navParams.get('projectId');
    this.showModal = navParams.get('modal');
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      contributorRequests: false,
      claManagerRequests: false,
      managers: true
    };
    this.sort = {
      date: {
        arrayProp: 'date_modified',
        sortType: 'date',
        sort: null
      }
    };

    this.company = new ClaCompanyModel();
    this.cclaSignature = new ClaSignatureModel();
  }

  ngOnInit() {
    this.getProject();
    this.getCompany();
    this.listPendingContributorRequests();
    this.listPendingCLAManagerRequests();
  }

  getProject() {
    this.loading.projects = true;
    this.claService.getProject(this.projectId).subscribe((response) => {
      this.loading.projects = false;
      this.project = response;
      this.projectName = this.project.project_name;
      this.getProjectSignatures();
    });
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe((response) => {
      this.company = response;
      this.getManager(this.company.company_manager_id);
    });
  }

  getCLAManagers() {
    this.loading.managers = true;
    this.claService.getCLAManagers(this.cclaSignature.signatureID).subscribe((response) => {
      this.loading.managers = false;
      if (response.errors != null) {
        this.managers = [];
      } else {
        this.managers = response;
        this.canEdit = this.hasEditAccess();
      }
    });
  }

  hasEditAccess() {
    for (let i = 0; i < this.managers.length; i++) {
      if (this.managers[i].email === this.userEmail) {
        return true;
      }
    }
    return false;
  }

  getProjectSignatures() {
    // Get CCLA Company Signatures - should just be one
    this.loading.signatures = true;
    this.claService.getCompanyProjectSignatures(this.companyId, this.projectId).subscribe(
      (response) => {
        this.loading.signatures = false;
        if (response.signatures) {
          let cclaSignatures = response.signatures.filter((sig) => sig.signatureType === 'ccla');
          if (cclaSignatures.length) {
            this.cclaSignature = cclaSignatures[0];

            // Sort the values
            if (this.cclaSignature.domainWhitelist) {
              const sortedList: string[] = this.cclaSignature.domainWhitelist.sort((a, b) => {
                return a.trim().localeCompare(b.trim());
              });
              // Remove duplicates - set doesn't allow dups
              this.cclaSignature.domainWhitelist = Array.from(new Set(sortedList));
            }
            if (this.cclaSignature.emailWhitelist) {
              const sortedList: string[] = this.cclaSignature.emailWhitelist.sort((a, b) => {
                return a.trim().localeCompare(b.trim());
              });
              // Remove duplicates - set doesn't allow dups
              this.cclaSignature.emailWhitelist = Array.from(new Set(sortedList));
            }
            if (this.cclaSignature.githubWhitelist) {
              const sortedList: string[] = this.cclaSignature.githubWhitelist.sort((a, b) => {
                return a.trim().localeCompare(b.trim());
              });
              // Remove duplicates - set doesn't allow dups
              this.cclaSignature.githubWhitelist = Array.from(new Set(sortedList));
            }
            if (this.cclaSignature.githubOrgWhitelist) {
              const sortedList: string[] = this.cclaSignature.githubOrgWhitelist.sort((a, b) => {
                return a.trim().localeCompare(b.trim());
              });
              // Remove duplicates - set doesn't allow dups
              this.cclaSignature.githubOrgWhitelist = Array.from(new Set(sortedList));
            }
            this.getCLAManagers();
            //this.getGitHubOrgWhitelist();
          }
        }
      },
      (exception) => {
        this.loading.signatures = false;
        console.log(exception);
      }
    );

    // Get CCLA Employee Signatures
    this.loading.acknowledgements = true;
    this.claService.getEmployeeProjectSignatures(this.companyId, this.projectId).subscribe(
      (response) => {
        this.loading.acknowledgements = false;
        if (response.signatures) {
          const sigs = response.signatures || [];
          this.employeeSignatures = sigs.sort((a, b) => {
            if (a.userName != null && b.userName != null) {
              return a.userName.trim().localeCompare(b.userName.trim());
            } else {
              return 0;
            }
          });
        }
      },
      (exception) => {
        this.loading.acknowledgements = false;
        console.log(exception);
      }
    );
  }

  getManager(userId) {
    this.claService.getUser(userId).subscribe((response) => {
      this.manager = response;
    });
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'domain',
      projectName: this.project.project_name,
      companyName: this.company.company_name,
      projectId: this.cclaSignature.projectID,
      signatureId: this.cclaSignature.signatureID,
      whitelist: this.cclaSignature.domainWhitelist
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'email',
      projectName: this.project.project_name,
      companyName: this.company.company_name,
      projectId: this.cclaSignature.projectID,
      companyId: this.companyId,
      signatureId: this.cclaSignature.signatureID,
      whitelist: this.cclaSignature.emailWhitelist
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistGithubModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'github',
      projectName: this.project.project_name,
      companyName: this.company.company_name,
      projectId: this.cclaSignature.projectID,
      signatureId: this.cclaSignature.signatureID,
      whitelist: this.cclaSignature.githubWhitelist
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistGithubOrgModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'githubOrg',
      projectName: this.project.project_name,
      companyName: this.company.company_name,
      projectId: this.cclaSignature.projectID,
      signatureId: this.cclaSignature.signatureID,
      whitelist: this.cclaSignature.githubOrgWhitelist
    });
    modal.onDidDismiss((data) => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  sortMembers(prop) {
    this.sortService.toggleSort(this.sort, prop, this.employeeSignatures);
  }

  deleteManagerConfirmation(payload: ClaManager) {
    let alert = this.alertCtrl.create({
      subTitle: `Delete Manager`,
      message: `Are you sure you want to delete this Manager?`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Delete',
          handler: () => {
            this.deleteManager(payload);
          }
        }
      ]
    });

    alert.present();
  }

  deleteManager(payload) {
    this.claService.deleteCLAManager(this.cclaSignature.signatureID, payload.lfid).subscribe(() =>
      this.getCLAManagers()
    );
  }

  openManagerModal() {
    let modal = this.modalCtrl.create('AddManagerModal', {
      signatureId: this.cclaSignature.signatureID
    });
    modal.onDidDismiss((data) => {
      if (data) {
        this.getCLAManagers();
      }
    });
    modal.present();
  }

  listPendingContributorRequests() {
    this.claService.getProjectWhitelistRequest(this.companyId, this.projectId, null).subscribe((request) => {
      if (request.list.length == 0) {
        this.noPendingContributorRequests = true;
      } else {
        this.allContributorRequests = request.list;
        this.pendingContributorRequests = request.list.filter((req) => req.requestStatus === 'pending');
        this.approvedContributorRequests = request.list.filter((req) => req.requestStatus === 'approved');
        this.rejectedContributorRequests = request.list.filter((req) => req.requestStatus === 'rejected');
        this.noPendingContributorRequests = false;
      }
      this.loading.contributorRequests = false;
    }, (error) => {
      console.log(`error loading contributor requests: ${error}`);
      this.loading.contributorRequests = false;
    })
  }

  listPendingCLAManagerRequests() {
    console.log('loading CLA Manager Requests');
    this.claService.getCLAManagerRequests(this.companyId, this.projectId).subscribe((response) => {
      //console.log('CLA Manager Requests:');
      //console.log(response);
      if (response.requests == null || response.requests.length == 0) {
        console.log('No Pending CLA Manager Requests')
        this.noPendingCLAManagerRequests = true;
      } else {
        this.allCLAManagerRequests = response.requests;
        this.pendingCLAManagerRequests = response.requests.filter((req) => req.status === 'pending');
        this.approvedCLAManagerRequests = response.requests.filter((req) => req.status === 'approved');
        this.deniedCLAManagerRequests = response.requests.filter((req) => req.status === 'denied');
        this.noPendingCLAManagerRequests = false;
      }
      this.loading.claManagerRequests = false;
    }, (error) => {
      console.log(`error loading cla manager requests: ${error}`);
      this.loading.claManagerRequests = false;
    })

  }

  accept(requestID) {
    let alert = this.alertCtrl.create({
      subTitle: `Approve Request - Confirmation`,
      message: 'This will dismiss this pending request and send the contributor ' +
        'an email confirming that they have access. <br/><br/>Please confirm that you ' +
        'have updated the Approved Lists by adding this user by email, domain, ' +
        'GitHub username or Organization.<br/><br/>Are you sure?',
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
            this.claService.approveCclaWhitelistRequest(this.companyId, this.projectId, requestID)
              .subscribe(
                (res) => {
                  this.listPendingContributorRequests();
                },
                (error) => console.log(error));
          }
        }
      ]
    });
    alert.present();
  }

  decline(requestID) {
    let alert = this.alertCtrl.create({
      subTitle: 'Decline Request - Confirmation',
      message: 'This will dismiss this pending request and send the contributor ' +
        'an email confirming that their request was denied.' +
        '<br/><br/>Are you sure you want to decline this request?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Decline Request',
          handler: () => {
            this.claService.rejectCclaWhitelistRequest(this.companyId, this.projectId, requestID)
              .subscribe(
                (res) => {
                  this.listPendingContributorRequests();
                },
                (error) => console.log(error));
          }
        }
      ]
    });
    alert.present();
    alert.present();
  }

  approveCLAManagerRequest(companyID: string, companyName: string, projectID, projectName: string, requestID: string) {
    let alert = this.alertCtrl.create({
      subTitle: `Approve Request - Confirmation`,
      message: `This will dismiss this pending CLA Manager request and send all the CLA Managers for ${companyName} `+
        `for project ${projectName} ` +
        'an email confirming that this user will now have CLA Manager authority. '+
        'Additionally, the new CLA Manager will also receive an email with instructions '+
        'on how to login, view and manage the approval lists.'+
        '<br/><br/>Are you sure you want to approve this request?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Approve Request',
          handler: () => {
            this.claService.approveCLAManagerRequest(companyID, projectID, requestID)
              .subscribe(
                (res) => {
                  this.listPendingCLAManagerRequests();
                },
                (error) => console.log(error));
          }
        }
      ]
    });
    alert.present();
  }

  denyCLAManagerRequest(companyID: string, companyName: string, projectID, projectName: string, requestID: string) {
    let alert = this.alertCtrl.create({
      subTitle: 'Decline Request - Confirmation',
      message: 'This will dismiss this CLA Manager request and send the current CLA Managers an email indicating that this request was denied.' +
        '<br/><br/>Are you sure you want to decline this request?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {
          }
        },
        {
          text: 'Decline Request',
          handler: () => {
            this.claService.denyCLAManagerRequest(companyID, projectID, requestID)
              .subscribe(
                (res) => {
                  this.listPendingCLAManagerRequests();
                },
                (error) => console.log(error));
          }
        }
      ]
    });
    alert.present();
    alert.present();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  back() {
    this.navCtrl.setRoot('CompanyPage', {
      companyId: this.companyId
    });
  }

  showRequestType(requestType: string) {
    // TODO: implement filter to show which items should be selected
  }
}
