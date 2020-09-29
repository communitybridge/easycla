// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { AppSettings } from 'src/app/config/app-settings';
import { EnvConfig } from 'src/app/config/cla-env-utils';
import { AuthService } from 'src/app/core/services/auth.service';

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
  styleUrls: ['./page-not-found.component.scss']
})
export class PageNotFoundComponent implements OnInit {
  message: string;

  constructor(
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    const consoleType: string = JSON.parse(this.authService.getItem(AppSettings.CONSOLE_TYPE));
    if (consoleType) {
      this.message = 'You are being redirected to the ' + (consoleType === AppSettings.PROJECT_CONSOLE_LINK ? 'Project' : 'Corporate') + ' contributor console.';
      const url = EnvConfig.default[consoleType] + '?idToken=' + this.authService.getIdToken();
      this.authService.removeItem(AppSettings.CONSOLE_TYPE);
      window.open(url, '_self');
    } else {
      this.message = 'The page you are looking for was not found.';
    }
  }

}
