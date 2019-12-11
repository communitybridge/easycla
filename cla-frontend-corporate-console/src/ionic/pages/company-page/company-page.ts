// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavController, NavParams } from 'ionic-angular';
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
  rows: any[];

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
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
    this.projects = [];

    this.data = {};
    this.columns = [{ prop: 'ProjectName' }, { prop: 'ProjectManagers' }];
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
    //console.log('Loading company signatures...');
    this.loading.companySignatures = true;
    this.loading.projects = true;

    // Clear out our projects and table models
    this.projects = [];
    this.rows = [];

    this.claService.getCompanySignatures(this.companyId).subscribe(
      response => {
        //console.log('Company signatures:');
        //console.log(response);
        if (response.resultCount > 0) {
          //console.log('Filtering Company signatures...');
          this.companySignatures = response.signatures.filter(signature => signature.signatureSigned === true);
          //console.log('Filtered Company signatures: ' + this.companySignatures.length);
          //console.log('Loading projects...');
          for (let signature of this.companySignatures) {
            this.getProject(signature.projectID);
          }
        }
        this.loading.companySignatures = false;
        this.loading.projects = false;
      },
      exception => {
        this.loading.companySignatures = false;
        this.loading.projects = false;
      }
    );
  }

  getProject(projectId) {
    //console.log('Loading project: ' + projectId);
    this.claService.getProject(projectId).subscribe(response => {
      //console.log('Loaded project: ');
      //console.log(response);
      this.projects.push(response);

      this.loading.projects = false;
      this.rows = this.mapProjects(this.projects);
    });
  }

  mapProjects(projects) {
    let rows = [];
    for (let project of projects) {
      rows.push({
        ProjectID: project.project_id,
        ProjectName: project.project_name,
        ProjectManagers: project.project_acl
      });
    }

    return rows;
  }

  onSelect(event) {
    this.openProjectPage(event.selected[0].ProjectID);
  }

  openProjectPage(projectId) {
    this.navCtrl.push('ProjectPage', {
      companyId: this.companyId,
      projectId: projectId
    });
  }

  openCompanyModal() {
    let modal = this.modalCtrl.create('AddCompanyModal', {
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create('WhitelistModal', {
      type: 'email',
      company: this.company
    });
    modal.onDidDismiss(data => {
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
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openProjectsCclaSelectModal() {
    let modal = this.modalCtrl.create('ProjectsCclaSelectModal', {
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
    });
  }

  declineCompanyInvite(invite) {
    let data = {
      inviteId: invite.inviteId,
      userLFID: invite.userLFID
    };
    this.claService.declineCompanyInvite(this.companyId, data).subscribe(response => {
      this.getInvites();
    });
  }
}
