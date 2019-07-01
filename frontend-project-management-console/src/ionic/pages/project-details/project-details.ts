// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, Content } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UrlValidator } from  '../../validators/url';
import { forbiddenValidator } from '../../validators/forbidden';
import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { ProjectModel } from '../../models/project-model';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project-details/:projectId'
// })
@Component({
  selector: 'project-details',
  templateUrl: 'project-details.html',
})
export class ProjectDetailsPage {

  keysGetter;
  projectId: string;

  project = new ProjectModel();
  projectStatuses: any;
  projectCategories: any;
  projectSectors: any;

  membershipsCount: number;

  editProject: any;
  loading: any;
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  @ViewChild(Content) content: Content;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    private formBuilder: FormBuilder,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');
    this.getDefaults();
    this.form = formBuilder.group({
      name:[this.project.name, Validators.compose([Validators.required])],
      startDate:[this.project.startDate],
      status:[this.project.status],
      category:[this.project.category],
      sector:[this.project.sector, Validators.compose([forbiddenValidator(/INVALID/i)])],
      url:[this.project.url, Validators.compose([UrlValidator.isValid])],
      addressThoroughfare:[this.project.address.address.thoroughfare],
      addressPostalCode:[this.project.address.address.postalCode],
      addressLocalityName:[this.project.address.address.localityName],
      addressAdministrativeArea:[this.project.address.address.administrativeArea],
      addressCountry:[this.project.address.address.country],
      description:[this.project.description],
    });
  }

  ngOnInit() {
    this.getProject(this.projectId);
    this.getProjectStatuses();
    this.getProjectCategories();
    this.getProjectSectors();
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if (response) {
        this.project = response;
        this.loading.project = false;
        if(this.project.config.logoRef) { this.project.config.logoRef += "?" + new Date().getTime(); }
        this.form.patchValue({
          name:this.project.name,
          startDate:this.project.startDate,
          status:this.project.status,
          category:this.project.category,
          sector:this.project.sector,
          url:this.project.url,
          addressThoroughfare:this.project.address.address.thoroughfare,
          addressPostalCode:this.project.address.address.postalCode,
          addressLocalityName:this.project.address.address.localityName,
          addressAdministrativeArea:this.project.address.address.administrativeArea,
          addressCountry:this.project.address.address.country,
          description:this.project.description,
        });
      }
    });
  }

  getProjectStatuses() {
    this.cincoService.getProjectStatuses().subscribe(response => {
      this.projectStatuses = response;
    });
  }

  getProjectCategories() {
    this.cincoService.getProjectCategories().subscribe(response => {
      this.projectCategories = response;
    });
  }

  getProjectSectors() {
    this.cincoService.getProjectSectors().subscribe(response => {
      this.projectSectors = response;
    });
  }

  submitEditProject() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.content.scrollToTop();
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let address = {
      address: {
        thoroughfare: this.form.value.addressThoroughfare,
        postalCode: this.form.value.addressPostalCode,
        localityName: this.form.value.addressLocalityName,
        administrativeArea: this.form.value.addressAdministrativeArea,
        country: this.form.value.addressCountry,
      },
      type: "BILLING",
    };

    let sector = (this.form.value.sector || undefined);
    this.editProject = {
      name: this.form.value.name,
      description: this.form.value.description,
      url: this.form.value.url,
      address: address,
      sector: sector,
      status: this.form.value.status,
      category: this.form.value.category,
      startDate: this.form.value.startDate,
    };

    this.cincoService.editProject(this.projectId, this.editProject).subscribe(
      response => {
        this.currentlySubmitting = false;
        this.navCtrl.setRoot('ProjectPage', {
          projectId: this.projectId
        });
      },
      error => {
        console.log('Project Save Error: ' + error);
        // allow recovery from error
        this.currentlySubmitting = false;
      }
    );
  }

  cancelEditProject() {
    this.navCtrl.setRoot('ProjectPage', {
      projectId: this.projectId
    });
  }

  openAssetManagementModal() {
    let modal = this.modalCtrl.create('AssetManagementModal', {
      projectId: this.projectId,
      projectName: this.project.name
    });
    modal.onDidDismiss(newlogoRef => {
      if(newlogoRef){
        this.project.config.logoRef = newlogoRef;
      }
    });
    modal.present();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    this.editProject = {};
    this.loading = {
      project: true,
    };
    this.projectStatuses = {};
    this.projectCategories = {};
    this.projectSectors = {};
    this.project = {
      id: "",
      name: "",
      description: "",
      managers: "",
      members: "",
      status: "",
      category: "",
      sector: "",
      url: "",
      logoRef: "",
      startDate: "",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: {
        address: {
          administrativeArea: "",
          country: "",
          localityName: "",
          postalCode: "",
          thoroughfare: ""
        },
        type: ""
      },
      config: {
        logoRef: ""
      }
    };
  }

}
