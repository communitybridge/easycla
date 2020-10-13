// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Events, NavController } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';

@Component({
  selector: 'cla-header',
  templateUrl: 'cla-header.html'
})
export class ClaHeader {
  @Input() title = '';
  @Input() hasShowBackBtn = false;
  @Output() onToggle: EventEmitter<any> = new EventEmitter<any>();
  hasEnabledLFXHeader = EnvConfig['lfx-header-enabled'] === "true" ? true : false;
  hasExpanded: boolean = true;

  constructor(
    public navCtrl: NavController,
  ) { }

  onToggled() {
    this.hasExpanded = !this.hasExpanded;
    this.onToggle.emit(this.hasExpanded);
  }


  backToProjects() {
    this.navCtrl.pop();
  }

}

