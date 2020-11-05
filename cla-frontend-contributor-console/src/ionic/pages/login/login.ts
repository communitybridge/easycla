// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { AUTH_ROUTE } from '../../services/auth.utils';

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
  expanded: boolean = true;

  constructor(
    public authService: AuthService
  ) { }

  login() {
    this.authService.login(AUTH_ROUTE);
  }

  onClickToggle(hasExpanded) {
    this.expanded = hasExpanded;
  }
}
