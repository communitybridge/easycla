// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { ClaService } from '../../services/cla.service';

@Component({
  selector: 'cla-footer',
  templateUrl: 'cla-footer.html'
})
export class ClaFooter {
  version: any;
  releaseDate: any;
  constructor(
    public claService: ClaService,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.getReleaseVersion()
  }

  getReleaseVersion() {
    this.claService.getReleaseVersion().subscribe((data) => {
      this.version = data.version;
      this.releaseDate = data.buildDate;
    })
  }

  
}
