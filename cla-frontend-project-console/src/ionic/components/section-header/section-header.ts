// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';
import { NavController, ModalController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { ClaService } from '../../services/cla.service';

@Component({
  selector: 'section-header',
  templateUrl: 'section-header.html'
})
export class SectionHeaderComponent {
  @Input('projectId') private projectId: string;

  project: any;

  // projectSectors: any;

  constructor(
    private navCtrl: NavController,
    private cincoService: CincoService,
    private claService: ClaService,
    public modalCtrl: ModalController
  ) {
    this.getDefaults();
  }

  ngOnInit() {
    // this.getProjectSectors();
    this.getProject(this.projectId);
  }

  getDefaults() {
    // this.projectSectors = {};
    this.project = {
      id: '',
      name: 'Project',
      logoRef: '',
      description: ''
    };
  }

  // viewProjectDetails(projectId) {
  //   this.navCtrl.push("ProjectDetailsPage", {
  //     projectId: projectId
  //   });
  // }

  getProject(projectId) {
    // let getMembers = true;
    this.claService.getProjectFromSFDC(projectId).subscribe(response => {
      if (response) {
        this.project = response;
      }
    });
  }

  // getProjectSectors() {
  //   this.cincoService.getProjectSectors().subscribe(response => {
  //     this.projectSectors = response;
  //   });
  // }

  // openProjectUserManagementModal() {
  //   let modal = this.modalCtrl.create("ProjectUserManagementModal", {
  //     projectId: this.projectId,
  //     projectName: this.project.name
  //   });
  //   modal.present();
  // }

  // openAssetManagementModal() {
  //   let modal = this.modalCtrl.create("AssetManagementModal", {
  //     projectId: this.projectId,
  //     projectName: this.project.name
  //   });
  //   modal.present();
  // }
}
