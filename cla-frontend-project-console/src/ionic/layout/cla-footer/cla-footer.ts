import { Component } from '@angular/core';
import { ClaService } from '../../services/cla.service';
import { generalConstants } from '../../constants/general';

@Component({
  selector: 'cla-footer',
  templateUrl: 'cla-footer.html'
})
export class ClaFooter {
  version: any;
  releaseDate: any;
  helpURL: string = generalConstants.getHelpURL;
  acceptableUsePolicyURL: string = generalConstants.acceptableUsePolicyURL;
  serviceSpecificTermsURL: string = generalConstants.serviceSpecificTermsURL;
  platformUseAgreementURL: string = generalConstants.platformUseAgreementURL;
  privacyPolicyURL: string = generalConstants.privacyPolicyURL;

  constructor(
    public claService: ClaService,
  ) {
    this.getReleaseVersion();
  }

  getReleaseVersion() {
    this.claService.getReleaseVersion().subscribe((data) => {
      this.version = data.version;
      this.releaseDate = data.buildDate;
    })
  }

}


