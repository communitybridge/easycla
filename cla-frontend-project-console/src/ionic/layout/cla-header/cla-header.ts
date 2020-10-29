// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Events } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';

@Component({
  selector: 'cla-header',
  templateUrl: 'cla-header.html'
})
export class ClaHeader {
  @Input() title = '';
  @Input() hasShowBackBtn = false;
  @Output() onToggle: EventEmitter<any> = new EventEmitter<any>();
  hasExpanded: boolean = true;

  constructor(
    public events: Events,
  ) { }

  onToggled() {
    this.hasExpanded = !this.hasExpanded;
    this.onToggle.emit(this.hasExpanded);
  }


  backToProjects() {
    this.events.publish('nav:allProjects');
  }

}

