// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';

import { NavController, IonicPage, ModalController } from 'ionic-angular';

import { CincoService } from '../../services/cinco.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser', 'isAdmin']
})
// @IonicPage({
//   segment: 'console-users'
// })
@Component({
  selector: 'console-users',
  templateUrl: 'console-users.html'
})
export class ConsoleUsersPage {
  users: any;
  userRoles: any;
  loading: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    public rolesService: RolesService
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      users: true
    };
    this.users = [];
    this.userRoles = {};
  }

  ngOnInit() {
    this.getUserRoles();
    // this.getAllUsers();
  }

  getUserRoles() {
    this.cincoService.getUserRoles().subscribe(response => {
      if (response) {
        this.userRoles = response;
      }
    });
  }

  // TODO: Get all Users is not avaiable anymore in CINCO.
  // We need logic for search users instead now.
  // getAllUsers() {
  //   this.cincoService.getAllUsers().subscribe(response => {
  //     if(response) {
  //       this.users = response;
  //       this.loading.users = false;
  //     }
  //   });
  // }

  userSelected(user) {
    let modal = this.modalCtrl.create('ConsoleUserUpdateModal', {
      user: user
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      // this.getAllUsers();
    });
    modal.present();
  }

  addNewUser() {
    let modal = this.modalCtrl.create('ConsoleUserUpdateModal', {});
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      // this.getAllUsers();
    });
    modal.present();
  }
}
