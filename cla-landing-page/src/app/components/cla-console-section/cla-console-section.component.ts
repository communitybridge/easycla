// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input, OnInit } from '@angular/core';
import { AppSettings } from 'src/app/config/app-settings';
import { StorageService } from 'src/app/core/services/storage.service';
import { EnvConfig } from 'src/app/config/cla-env-utils';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-cla-console-section',
  templateUrl: './cla-console-section.component.html',
  styleUrls: ['./cla-console-section.component.scss']
})
export class ClaConsoleSectionComponent implements OnInit {
  @Input() consoleMetadata: any;
  version: string;
  links: any[];

  constructor(
    private storageService: StorageService,
    private router: ActivatedRoute
  ) { }

  ngOnInit() {
    this.version = this.router.snapshot.queryParamMap.get('version');
    const element: any = document.getElementById('lfx-header');
    let projectConsoleUrl = EnvConfig.default[AppSettings.PROJECT_CONSOLE_LINK] + '#/login';
    let corporateConsoleUrl = EnvConfig.default[AppSettings.CORPORATE_CONSOLE_LINK] + '#/login'
    if (this.version === '2') {
      // Set redirect URL to new V2 console.
      projectConsoleUrl = EnvConfig.default[AppSettings.PROJECT_CONSOLE_LINK_V2];
      corporateConsoleUrl = EnvConfig.default[AppSettings.CORPORATE_CONSOLE_LINK_V2];
    }
    this.links = [
      {
        title: 'Project Login',
        url: projectConsoleUrl
      },
      {
        title: 'CLA Manager Login',
        url: corporateConsoleUrl
      },
      {
        title: 'Developer',
        url: AppSettings.CONTRIBUTORS_LEARN_MORE
      }
    ];
    element.links = this.links;
  }

  onClickProceed(type: string) {
    let url = '';
    this.storageService.setItem('type', type);
    if (this.version === '2') {
      // Set redirect URL to new V2 console.
      const envKey = (type === 'Projects') ? AppSettings.PROJECT_CONSOLE_LINK_V2 : AppSettings.CORPORATE_CONSOLE_LINK_V2;
      url = EnvConfig.default[envKey];
    } else {
      // Set redirect URL to new V1 console.
      const envKey = (type === 'Projects') ? AppSettings.PROJECT_CONSOLE_LINK : AppSettings.CORPORATE_CONSOLE_LINK;
      url = EnvConfig.default[envKey] + '#/login';
    }
    window.open(url, '_self');
  }

  onClickRequestAccess() {
    window.open(AppSettings.REQUEST_ACCESS_LINK, '_blank');
  }

  onClickLearnMore() {
    window.open(AppSettings.CONTRIBUTORS_LEARN_MORE, '_blank');
  }
}
