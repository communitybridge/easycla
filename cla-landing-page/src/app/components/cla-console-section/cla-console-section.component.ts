// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input, OnInit } from '@angular/core';
import { AuthService } from 'src/app/core/services/auth.service';
import { EnvConfig } from 'src/app/config/cla-env-utils';
import { AppSettings } from 'src/app/config/app-settings';


@Component({
  selector: 'app-cla-console-section',
  templateUrl: './cla-console-section.component.html',
  styleUrls: ['./cla-console-section.component.scss']
})
export class ClaConsoleSectionComponent implements OnInit {
  @Input() consoleMetadata: any;

  constructor(
    private authService: AuthService
  ) { }

  ngOnInit(): void {
  }


  onClickProceed(type: string) {
    const consoleType = type === 'Projects' ? AppSettings.PROJECT_CONSOLE_LINK : AppSettings.CORPORATE_CONSOLE_LINK;
    const url = EnvConfig.default[consoleType];
    window.open(url, '_self');
  }

  onClickRequestAccess() {
    window.open(AppSettings.REQUEST_ACCESS_LINK, '_blank');
  }

  onClickLearnMore() {
    window.open(AppSettings.LEARN_MORE, '_blank');
  }

}
