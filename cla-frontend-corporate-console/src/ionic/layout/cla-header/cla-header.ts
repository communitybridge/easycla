// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'cla-header',
  templateUrl: 'cla-header.html'
})
export class ClaHeader {
  @Input() title = '';
  @Output() onToggle: EventEmitter<any> = new EventEmitter<any>();

  hasExpanded: boolean = true;

  onToggled() {
    this.hasExpanded = !this.hasExpanded;
    this.onToggle.emit(this.hasExpanded);
  }
}

