// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'lfx-header',
  templateUrl: 'lfx-header.html'
})

export class lfxHeader {
  @Input() expanded;
  @Output() toggled: EventEmitter<any> = new EventEmitter<any>();
}
