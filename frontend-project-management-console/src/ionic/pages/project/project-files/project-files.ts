// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project/:projectId/files'
// })
@Component({
  selector: 'project-files',
  templateUrl: 'project-files.html',
})
export class ProjectFilesPage {
  projectId: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {

  }


}
