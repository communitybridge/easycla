// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'lfx-header',
  templateUrl: './lfx-header.component.html',
  styleUrls: ['./lfx-header.component.scss']
})
export class LfxHeaderComponent implements OnInit {
  @Input() expanded: boolean;

  constructor() {
  }

  ngOnInit(): void {
  }
}
