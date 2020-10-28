// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { AUTH_ROUTE } from '../../services/auth.utils';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login',
  templateUrl: 'login.html'
})
export class LoginPage {
  canAccess: boolean;
  hasEnabledLFXHeader = EnvConfig['lfx-header-enabled'] === "true" ? true : false;
  
  constructor(
    public authService: AuthService
  ) { }

  login() {
    this.authService.login(AUTH_ROUTE);
  }

  onClickToggle(toggle) {

  }
}
