// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {AlertController, IonicPage, ModalController, NavController, NavParams} from "ionic-angular";
import {ClaService} from "../../services/cla.service";
import {ClaCompanyModel} from "../../models/cla-company";
import {ClaUserModel} from "../../models/cla-user";
import {ClaSignatureModel} from "../../models/cla-signature";
import {ClaManager} from "../../models/cla-manager";
import {SortService} from "../../services/sort.service";
import {RolesService} from "../../services/roles.service";
import {Restricted} from "../../decorators/restricted";
import {WhitelistModal} from "../../modals/whitelist-modal/whitelist-modal";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: "company/:companyId/project/:projectId/:modal"
})
@Component({
  selector: "project-page",
  templateUrl: "project-page.html"
})
export class ProjectPage {
  cclaSignature: ClaSignatureModel;
  employeeSignatures: ClaSignatureModel[];
  githubOrgWhitelist: any[] = [];
  githubEnabledWhitelist: any[] = [];
  loading: any;
  companyId: string;
  projectId: string;
  managers: ClaManager[];
  company: ClaCompanyModel;
  manager: ClaUserModel;
  showModal: any;

  project: any;
  users: any;

  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public alertCtrl: AlertController,
    public modalCtrl: ModalController,
    private rolesService: RolesService, // for @Restricted
    private sortService: SortService
  ) {
    this.companyId = navParams.get("companyId");
    this.projectId = navParams.get("projectId");
    this.showModal = navParams.get("modal");

    if (this.showModal === 'orgwhitelist') {
      this.openGithubOrgWhitelistModal();
    }

    this.getDefaults();
  }

  getDefaults() {
    this.loading = {};
    this.users = {};
    this.sort = {
      date: {
        arrayProp: "date_modified",
        sortType: "date",
        sort: null
      }
    };
    this.company = new ClaCompanyModel();
    this.cclaSignature = new ClaSignatureModel();
  }

  ngOnInit() {
    this.getProject();
    this.getCompany();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      this.company = response;
      this.getManager(this.company.company_manager_id);
    });
  }

  getProject() {
    this.claService.getProject(this.projectId).subscribe(response => {
      this.project = response;
      this.getProjectSignatures();
    });
  }

  getGitHubOrgWhitelist() {
    // console.log('loading GH Org Whitelist...');
    this.claService.getGithubOrganizationWhitelistEntries(this.cclaSignature.signature_id).subscribe(organizations => {
      // console.log('received GH Org Whitelist response: ');
      // console.log(organizations);
      this.githubOrgWhitelist = organizations;
    });
  }

  githubOrgWhitelistEnabled() {
    this.githubEnabledWhitelist = this.githubOrgWhitelist.filter((org) => org.selected);
  }

  getCLAManagers() {
    this.claService.getCLAManagers(this.cclaSignature.signature_id).subscribe(response => {
      this.managers = response;
    });
  }

  getProjectSignatures() {
    // get CCLA signatures
    this.claService
      .getCompanyProjectSignatures(this.companyId, this.projectId)
      .subscribe(response => {
        let cclaSignatures = response.filter(sig => sig.signature_type === 'ccla');
        if (cclaSignatures.length) {
          this.cclaSignature = cclaSignatures[0];
          this.getCLAManagers();
          this.getGitHubOrgWhitelist();
        }
      });

    // get employee signatures
    this.claService
      .getEmployeeProjectSignatures(this.companyId, this.projectId)
      .subscribe(response => {
        this.employeeSignatures = response;
        for (let signature of this.employeeSignatures) {
          this.getUser(signature.signature_reference_id);
        }
      });

  }

  getManager(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.manager = response;
    });
  }

  getUser(userId) {
    if (!this.users[userId]) {
      this.claService.getUser(userId).subscribe(response => {
        this.users[userId] = response;
      });
    }
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "email",
      signatureId: this.cclaSignature.signature_id,
      whitelist: this.cclaSignature.email_whitelist
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "domain",
      signatureId: this.cclaSignature.signature_id,
      whitelist: this.cclaSignature.domain_whitelist
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistGithubModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "github",
      signatureId: this.cclaSignature.signature_id,
      whitelist: this.cclaSignature.github_whitelist
    });
    modal.onDidDismiss(data => {
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
        }, {
          text: 'Delete',
          handler: () => {
            this.deleteManager(payload);
          }
        }
      ]
    });

    alert.present();
  }

  deleteManager(payload: ClaManager) {
    this.claService.deleteCLAManager(this.cclaSignature.signature_id, payload)
      .subscribe(() => this.getCLAManagers())
  }

  openManagerModal() {
    let modal = this.modalCtrl.create("AddManagerModal", {
      signatureId: this.cclaSignature.signature_id
    });
    modal.onDidDismiss(data => {
      if (data) {
        this.getCLAManagers();
      }
    });
    modal.present();
  }

  openGithubOrgWhitelistModal() {
    // Maybe this happens if we authorize GH and come back after the flow and we haven't fully loaded...
    if (this.cclaSignature == null) {
      this.claService.getCompanyProjectSignatures(this.companyId, this.projectId)
        .subscribe(response => {
          let cclaSignatures = response.filter(sig => sig.signature_type === 'ccla');
          if (cclaSignatures.length) {
            this.cclaSignature = cclaSignatures[0];
            this.getCLAManagers();
            this.getGitHubOrgWhitelist();

            // Ok to open the modal now that we have signatures loaded
            let modal = this.modalCtrl.create("GithubOrgWhitelistModal", {
              companyId: this.companyId,
              corporateClaId: this.projectId,
              signatureId: this.cclaSignature.signature_id
            });
            modal.onDidDismiss(data => {
              // Refresh the list
              this.getGitHubOrgWhitelist();
            });
            modal.present();
          }
        });
    } else {
      let modal = this.modalCtrl.create("GithubOrgWhitelistModal", {
        companyId: this.companyId,
        corporateClaId: this.projectId,
        signatureId: this.cclaSignature.signature_id
      });
      modal.onDidDismiss(data => {
        // Refresh the list
        this.getGitHubOrgWhitelist();
      });
      modal.present();
    }
  }
}
