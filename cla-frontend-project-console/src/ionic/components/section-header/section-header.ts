// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';
import { ModalController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@Component({
  selector: 'section-header',
  templateUrl: 'section-header.html'
})
export class SectionHeaderComponent {
  @Input('projectId') private projectId: string;

  project: any = null;

  constructor(
    private claService: ClaService,
    public modalCtrl: ModalController
  ) { }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    this.claService.getProjectFromSFDC(projectId).subscribe((response) => {
      if (response) {
        this.project = response;
      }
    });
  }
}
