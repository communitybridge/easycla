// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {IonicPage, ModalController, NavController, NavParams} from "ionic-angular";
import {ClaService} from "../../services/cla.service";
import {ClaCompanyModel} from "../../models/cla-company";
import {ClaUserModel} from "../../models/cla-user";
import {ClaSignatureModel} from "../../models/cla-signature";
import {RolesService} from "../../services/roles.service";
import {Restricted} from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: "company/:companyId"
})
@Component({
  selector: "company-page",
  templateUrl: "company-page.html"
})
export class CompanyPage {
  companyId: string;
  company: ClaCompanyModel;
  manager: ClaUserModel;
  companySignatures: ClaSignatureModel[];
  projects: any;
  loading: any;
  invites: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private rolesService: RolesService // for @Restricted
  ) {
    this.companyId = navParams.get("companyId");
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companySignatures: true,
      invites: true
    };
    this.company = new ClaCompanyModel();
    this.projects = {};
  }

  ngOnInit() {
    this.getCompany();
    this.getCompanySignatures();
    this.getInvites();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      this.company = response;
      this.getUser(this.company.company_manager_id);
    });
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.manager = response;
    });
  }

  getCompanySignatures() {
    this.claService.getCompanySignatures(this.companyId).subscribe(response => {
      console.log(response);
      if (response.resultCount > 0) {
        this.companySignatures = response.signatures.filter(signature =>
          signature.signature_signed === true
        );
        for (let signature of this.companySignatures) {
          this.getProject(signature.signature_project_id);
        }
      }
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.projects[projectId] = response;
    });
  }

  openProjectPage(projectId) {
    this.navCtrl.push("ProjectPage", {
      companyId: this.companyId,
      projectId: projectId
    });
  }

  openCompanyModal() {
    let modal = this.modalCtrl.create("AddCompanyModal", {
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "email",
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "domain",
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openProjectsCclaSelectModal() {
    let modal = this.modalCtrl.create("ProjectsCclaSelectModal", {
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  getInvites() {
    this.claService.getPendingInvites(this.companyId).subscribe(response => {
      this.invites = response;
      this.loading.invites = false;
    });
  }

  acceptCompanyInvite(invite) {
    let data = {
      inviteId: invite.inviteId,
      userLFID: invite.userLFID
    };
    this.claService.acceptCompanyInvite(this.companyId, data).subscribe(response => {
      this.getInvites();
    })
  }

  declineCompanyInvite(invite) {
    let data = {
      inviteId: invite.inviteId,
      userLFID: invite.userLFID
    };
    this.claService.declineCompanyInvite(this.companyId, data).subscribe(response => {
      this.getInvites();
    })
  }
}
