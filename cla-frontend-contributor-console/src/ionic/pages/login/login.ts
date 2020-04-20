// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavParams } from 'ionic-angular';
import { RolesService } from '../../services/roles.service';
import { AuthService } from '../../services/auth.service';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login',
  templateUrl: 'login.html'
})
export class LoginPage {
  userRoles: any;
  canAccess: boolean;

  constructor(
    public navParams: NavParams,
    public rolesService: RolesService,
    public authService: AuthService
  ) {
    this.userRoles = this.rolesService.userRoles;
    this.rolesService.getUserRolesPromise().then((userRoles) => {
      this.userRoles = userRoles;
      this.canAccess = this.hasAccess();
    });
  }

  hasAccess() {
    return this.userRoles.isAuthenticated;
  }

  login() {
    this.authService.login();
  }
}
