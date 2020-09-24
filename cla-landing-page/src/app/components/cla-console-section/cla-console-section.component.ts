// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-cla-console-section',
  templateUrl: './cla-console-section.component.html',
  styleUrls: ['./cla-console-section.component.scss']
})
export class ClaConsoleSectionComponent implements OnInit {
  @Input() consoleMetadata: any;

  constructor() { }

  ngOnInit(): void {
  }


  onClickSignIn() {

  }

  onClickSignUp() {

  }

  onClickLearnMore() {

  }

}
