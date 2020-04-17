// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { AlertController, Events, IonicPage, ModalController, Nav, NavController, NavParams } from 'ionic-angular';
import { ClaService } from '../../../services/cla.service';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';
import { GithubOrganisationModel } from '../../../models/github-organisation-model';
import { PlatformLocation } from '@angular/common';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser']
})
@IonicPage({
  segment: 'project/:projectId/cla'
})
@Component({
  selector: 'project-cla',
  templateUrl: 'project-cla.html'
})
export class ProjectClaPage {
  loading: any;

  sfdcProjectId: string;
  githubOrganizations: GithubOrganisationModel[];

  claProjects: any = [];
  projectsByExternalId: any;
  alert;
  iclaUploadInfo: any;
  cclaUploadInfo: any;
  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public alertCtrl: AlertController,
    public claService: ClaService,
    public rolesService: RolesService,
    public events: Events,
    private location: PlatformLocation
  ) {
    this.sfdcProjectId = navParams.get('projectId');
    this.location.onPopState(() => {
      if (this.alert) {
        this.alert.dismiss();
      }
    });
    this.getDefaults();
    this.listenEvents();
  }

  listenEvents() {
    this.events.subscribe('reloadProjectCla', () => {
      this.getClaProjects();
    });
  }

  getDefaults() {
    this.loading = {
      claProjects: true,
      orgs: true
    };
    this.claProjects = [];
  }

  ngOnInit() {
    this.getClaProjects();
    this.getGithubOrganisation();
  }

  setLoadingOrganizationsSpinner(value) {
    this.loading = {
      orgs: value
    };
  }

  sortClaProjects(projects) {
    if (projects == null || projects.length == 0) {
      return projects;
    }

    return projects.sort((a, b) => {
      return a.projectName.trim().localeCompare(b.projectName.trim());
    });
  }

  sortGithubOrganisation(githubOrganizations) {
    if (githubOrganizations == null || githubOrganizations.length == 0) {
      return githubOrganizations;
    }

    return githubOrganizations.sort((a, b) => {
      return this.getOrganisationName(a.organizationName).trim().localeCompare(this.getOrganisationName(b.organizationName).trim());
    });
  }

  getClaProjects() {
    this.loading.claProjects = true;
    this.claService.getProjectsByExternalId(this.sfdcProjectId).subscribe((response) => {
      this.loading.claProjects = false;
      this.projectsByExternalId = response.projects;
      this.claProjects = this.sortClaProjects(response.projects);
    })
  }

  getGithubOrganisation() {
    this.loading.orgs = true;
    this.claService.getOrganizations(this.sfdcProjectId).subscribe((organizations) => {
      this.loading.orgs = false;
      this.githubOrganizations = this.sortGithubOrganisation(organizations.list);
    });
  }

  getOrganisationName(url) {
    return url.split('/').pop();
  }

  backToProjects() {
    this.events.publish('nav:allProjects');
  }

  openClaContractConfigModal(claProject) {
    let modal;
    // Edit CLA Group
    if (claProject) {
      modal = this.modalCtrl.create('ClaContractConfigModal', {
        claProject: claProject
      });
    } else {
      // Create CLA Group
      modal = this.modalCtrl.create('ClaContractConfigModal', {
        projectId: this.sfdcProjectId
      });
    }
    modal.onDidDismiss((data) => {
      if (data) {
        this.getClaProjects();
      }

    });
    modal.present();
  }

  goToSelectTemplatePage(projectId) {
    this.navCtrl.push('ProjectClaTemplatePage', {
      sfdcProjectId: this.sfdcProjectId,
      projectId: projectId
    });
  }

  openClaViewSignaturesModal(projectID: string, projectName: string) {
    let modal = this.modalCtrl.create(
      'ClaContractViewSignaturesModal',
      {
        claProjectId: projectID,
        claProjectName: projectName
      },
      {
        cssClass: 'medium'
      }
    );

    modal.present().catch((error) => {
      console.log('Error opening signatures modal view, error: ' + error);
    });
  }

  openClaViewCompaniesModal(projectID: string, projectName: string) {
    let modal = this.modalCtrl.create(
      'ClaContractViewCompaniesSignaturesModal',
      {
        claProjectId: projectID,
        claProjectName: projectName
      },
      {
        cssClass: 'medium'
      }
    );
    modal.present().catch((error) => {
      console.log('Error opening signatures modal view, error: ' + error);
    });
  }

  openClaContractVersionModal(claProjectId, documentType, documents) {
    let modal = this.modalCtrl.create('ClaContractVersionModal', {
      claProjectId: claProjectId,
      documentType: documentType,
      documents: documents
    });
    modal.present();
  }

  openClaOrganizationProviderModal(claProjectId) {
    let modal = this.modalCtrl.create('ClaOrganizationProviderModal', {
      claProjectId: claProjectId
    });
    modal.onDidDismiss((data) => {
      if (data) {
        this.openClaOrganizationAppModal();
      }
    });
    modal.present();
  }

  openClaConfigureGithubRepositoriesModal(claProjectId) {
    let modal = this.modalCtrl.create('ClaConfigureGithubRepositoriesModal', {
      claProjectId: claProjectId
    });
    modal.onDidDismiss((data) => {
      this.getClaProjects();
    });
    modal.present();
  }

  openClaGerritModal(projectId) {
    let modal = this.modalCtrl.create('ClaGerritModal', {
      projectId: projectId
    });
    modal.onDidDismiss((data) => {
      if (data) {
        this.getClaProjects();
      }
    });
    modal.present();
  }

  openClaOrganizationAppModal() {
    let modal = this.modalCtrl.create('ClaOrganizationAppModal', {});
    modal.onDidDismiss((data) => {
      this.getGithubOrganisation();
    });
    modal.present();
  }

  searchProjects(name: string, projects: any) {
    let found = false;

    if (projects) {
      projects.forEach((project) => {
        if (project.projectName.search(name) !== -1) {
          found = true;
        }
      });

    }

    return found;
  }

  deleteConfirmation(type, payload) {
    this.alert = this.alertCtrl.create({
      subTitle: `Delete ${type}`,
      message: `Are you sure you want to delete this ${type}?`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => { }
        },
        {
          text: 'Delete',
          handler: () => {
            switch (type) {
              case 'Github Organization':
                this.deleteClaGithubOrganization(payload);
                break;

              case 'Gerrit Instance':
                this.deleteGerritInstance(payload);
                break;
            }
          }
        }
      ]
    });

    this.alert.present();
  }


  deleteClaGroup(name: string, id: string) {
    this.alert = this.alertCtrl.create({
      subTitle: `Delete ${name}`,
      message: `Are you sure you want to delete this ${name}?`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => { }
        },
        {
          text: 'Delete',
          handler: () => {
            this.deleteClaProject(name, id);
          }
        }
      ]
    });

    this.alert.present();
  }

  deleteClaGroupError(name: string, id: string) {
    this.alert = this.alertCtrl.create({
      subTitle: `Delete ${name}`,
      message: `Could not delete ${name}. Please try again or contact support`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => { }
        },
        {
          text: 'Retry',
          handler: () => {
            this.deleteClaProject(name, id);
          }
        }
      ]
    });

    this.alert.present();
  }

  deleteClaGroupSuccess(name: string) {
    this.alert = this.alertCtrl.create({
      subTitle: `Delete ${name}`,
      message: `${name} was deleted successfully`,
      buttons: [
        {
          text: 'Close',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => { }
        }
      ]
    });

    this.alert.present();
  }


  deleteClaProject(name: string, id: string) {
    this.loading.claProject = true;
    this.claService.deleteClaProject(id).subscribe((res: any) => {
      if (res.status === 204) {
        this.getClaProjects();
        this.deleteClaGroupSuccess(name);
      }
      else {
        this.deleteClaGroupError(name, id)
        this.loading.claProject = false;
      }
    }, err => {
      this.deleteClaGroupError(name, id)
      this.loading.claProject = false;
    })
  }


  deleteClaGithubOrganization(organization) {
    this.claService.deleteGithubOrganization(this.sfdcProjectId, this.getOrganisationName(organization.organizationName)).subscribe((response) => {
      if (response.status === 200) {
        this.getGithubOrganisation();
      }
    });
  }

  deleteGerritInstance(gerrit) {
    this.claService.deleteGerritInstance(gerrit.gerrit_id).subscribe((response) => {
      this.getClaProjects();
    });
  }

  /**
   * Called if popover dismissed with data. Passes data to a callback function
   * @param  {object} popoverData should contain .callback and .callbackData
   */
  popoverResponse(popoverData) {
    let callback = popoverData.callback;
    if (this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  trimCharacter(text, length) {
    return text.length > length ? text.substring(0, length) + '...' : text;
  }

}