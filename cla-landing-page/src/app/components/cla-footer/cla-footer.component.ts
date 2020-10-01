// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-cla-footer',
  templateUrl: './cla-footer.component.html',
  styleUrls: ['./cla-footer.component.scss']
})
export class ClaFooterComponent implements OnInit {

  COMMUNITY_PEOPLE_GUIDE = 'https://communitybridge.org/guide/';
  PRIVACY_POLICY = 'https://www.linuxfoundation.org/privacy/';
  TRADEMARK_USEAGE = 'https://www.linuxfoundation.org/trademark-usage/';
  ACCEPTABLE_USER_POLICY = environment.ACCEPTABLE_USER_POLICY;
  SERVICE_SPECIFIC_TERM = environment.SERVICE_SPECIFIC_TERM
  PLATEFORM_USER_AGREEMENT = environment.PLATEFORM_USER_AGREEMENT;

  constructor() { }

  ngOnInit(): void {
  }

}
