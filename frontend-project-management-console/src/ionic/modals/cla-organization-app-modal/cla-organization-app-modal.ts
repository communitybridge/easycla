// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from "@angular/core";
import {
  NavController,
  ModalController,
  NavParams,
  ViewController,
  IonicPage, Events
} from "ionic-angular";
import { PopoverController } from "ionic-angular";
import { EnvConfig } from "../../services/cla.env.utils";
import {ClaService} from "../../services/cla.service";

@IonicPage({
  segment: "cla-organization-app-modal"
})
@Component({
  selector: "cla-organization-app-modal",
  templateUrl: "cla-organization-app-modal.html"
})
export class ClaOrganizationAppModal {
  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
    public events: Events
  ) {
    this.getDefaults();

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {}

  getDefaults() {}

  // TODO: Do we want a call to cla that polls for the installation status?
  // UH YEA?

  openAppPage() {
    this.viewCtrl.dismiss();
    window.open(EnvConfig['gh-app-public-link'], '_blank');
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
