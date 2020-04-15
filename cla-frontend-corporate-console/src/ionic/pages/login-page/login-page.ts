// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login-page',
  templateUrl: 'login-page.html'
})
export class LoginPage {
  constructor(
    private authService: AuthService
  ) { }

  login() {
    this.authService.login();
  }
}
