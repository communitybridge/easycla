// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ElementRef, ViewChild } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla/project/:projectId/new-company'
})
@Component({
  selector: 'cla-new-company-modal',
  templateUrl: 'cla-new-company-modal.html'
})
export class ClaNewCompanyModal {
  projectId: string;
  repositoryId: string;
  userId: string;

  consoleLink: string;

  @ViewChild('textArea') textArea: ElementRef;

  constructor(public navCtrl: NavController, public navParams: NavParams, public viewCtrl: ViewController) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.getDefaults();
  }

  getDefaults() {
    this.consoleLink = EnvConfig['corp-console-link'];
  }

  ngOnInit() {}

  dismiss() {
    this.viewCtrl.dismiss();
  }

  openConsoleLink() {
    window.open(this.consoleLink, '_blank');
  }

  copyText() {
    let copyTextarea = this.textArea.nativeElement;
    copyTextarea.select();

    try {
      var successful = document.execCommand('copy');
      var msg = successful ? 'successful' : 'unsuccessful';
      console.log('Copying text command was ' + msg);
    } catch (err) {
      console.log('Exception thrown attempting to copy to clipboard');
    }
  }
}
