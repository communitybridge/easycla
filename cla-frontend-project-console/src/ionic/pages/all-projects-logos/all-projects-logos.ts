// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';

import { NavController, IonicPage, ModalController } from 'ionic-angular';

import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';

import { ProjectModel } from '../../models/project-model';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'all-projects-logos'
// })
@Component({
  selector: 'all-projects-logos',
  templateUrl: 'all-projects-logos.html'
})
export class AllProjectsLogosPage {

  project = new ProjectModel();
  allProjects: any;
  projectLogos: any;
  sort: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
  ) {

  }

  async ngOnInit(){
    this.setSortProperties();
    this.getAllProjects();
  }

  getAllProjects() {
    this.cincoService.getAllProjects().subscribe(response => {
      if(response) {
        this.allProjects = response;
        for(let eachProject of this.allProjects) {
          this.getAllProjectsLogos(eachProject);
        }
      }
    });
  }

  getAllProjectsLogos(project) {
    this.cincoService.getProjectLogos(project.id).subscribe(response => {
      if(response) {
        this.projectLogos = response;
        project.logosCount = this.projectLogos.length
      }
    });
  }

  openAssetManagementModal(project) {
    let modal = this.modalCtrl.create('AssetManagementModal', {
      projectId: project.id,
      projectName: project.name
    });
    modal.onDidDismiss(newlogoRef => {
      if(newlogoRef){
        project.config.logoRef = newlogoRef;
      }
      this.getAllProjectsLogos(project);
    });
    modal.present();
  }


  setSortProperties(){
    this.sort = {
      project: {
        arrayProp: 'name',
        sortType: 'text',
        sort: null,
      },
      status: {
        arrayProp: 'status',
        sortType: 'text',
        sort: null,
      },
      logosCount: {
        arrayProp: 'logosCount',
        sortType: 'text',
        sort: null,
      }
    };
  }

  sortLogosTable(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.allProjects,
    );
  }

}
