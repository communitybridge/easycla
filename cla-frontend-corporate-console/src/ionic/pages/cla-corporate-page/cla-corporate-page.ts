// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {IonicPage, ModalController, NavController, NavParams} from "ionic-angular";
import {ClaService} from "../../services/cla.service";

@IonicPage({
  segment: "project/:projectId"
})
@Component({
  selector: "cla-corporate-page",
  templateUrl: "cla-corporate-page.html"
})
export class ClaCorporatePage {
  projectId: string;
  loading: any;
  company: any;
  project: any;
  signatureIntent: any;
  activeSignatures: boolean = true; // we assume true until otherwise
  signature: any;
  error: any = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService
  ) {
    this.getDefaults();
    this.projectId = navParams.get("projectId");
    this.company = navParams.get("company");
  }

  getDefaults() {
    this.project = {
      project_name: ""
    };
    this.signature = {
      sign_url: ""
    };
  }

  ngOnInit() {
    this.loading = true;
    this.getProject(this.projectId);
    this.postSignatureRequest();
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
    });
  }

  getReturnUrl() {
    return (
      window.location.protocol +
      "//" +
      window.location.host +
      "/#/company/" +
      this.company.company_id
    );
  }

  postSignatureRequest() {
    let signatureRequest = {
      project_id: this.projectId,
      company_id: this.company.company_id,
      // TODO: Switch this to intermediary loading screen as docusign postback has delay
      return_url: this.getReturnUrl()
    };

    this.claService
      .postCorporateSignatureRequest(signatureRequest)
      .subscribe(response => {
        this.loading = false;
        // returns {
        //   user_id:
        //   signature_id:
        //   project_id:
        //   sign_url: docusign.com/some-docusign-url
        // }
        if (response.errors) {
          this.error = response;
        }
        this.signature = response;
      });

  }

  openClaAgreement() {
    if (!this.signature.sign_url) {
      // Can't open agreement if we don't have a sign_url yet
      return;
    }
    window.open(this.signature.sign_url, "_blank");
  }
}
