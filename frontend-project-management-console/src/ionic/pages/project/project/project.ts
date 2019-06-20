// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from "@angular/core";
import {
  NavController,
  ModalController,
  NavParams,
  IonicPage
} from "ionic-angular";
import { ClaService } from "../../../services/cla.service";
import { CincoService } from "../../../services/cinco.service";
import { KeycloakService } from "../../../services/keycloak/keycloak.service";
import { SortService } from "../../../services/sort.service";
import { SFProjectModel } from "../../../models/sfdc-project-model";
import { RolesService } from "../../../services/roles.service";
import { Restricted } from "../../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated", "isPmcUser"]
})
@IonicPage({
  segment: "project/:projectId"
})
@Component({
  selector: "project",
  templateUrl: "project.html"
})
export class ProjectPage {
  selectedProject: any;
  projectId: string;

  project = new SFProjectModel();

  loading: any;
  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    private rolesService: RolesService,
    private claService: ClaService
  ) {
    this.projectId = navParams.get("projectId");
    this.getDefaults();
  }

  ngOnInit() {
    console.log("project id: " + this.projectId);
    // this.getSFDCProject(this.projectId);

    // this.cincoService
    //   .getEventsForProject(this.projectId)
    //   .subscribe(response => {
    //     if (response) {
    //       console.log(response);
    //     }
    //   });
  }

  getSFDCProject(projectId) {
    this.claService
      .getProjectFromSFDC(projectId)
      .subscribe(response => {
        if (response) {
          this.project = response;
          this.loading.project = false;
        }
      });
  }

  // getProject(projectId) {
  //   let getMembers = true;
  //   this.cincoService
  //     .getMockProject(projectId, getMembers)
  //     .subscribe(response => {
  //       if (response) {
  //         this.project = response;
  //         // // This is to refresh an image that have same URL
  //         // if (this.project.config.logoRef) {
  //         //   this.project.config.logoRef += "?" + new Date().getTime();
  //         // }
  //         this.loading.project = false;
  //       }
  //     });
  // }

  // memberSelected(event, memberId) {
  //   this.navCtrl.push("MemberPage", {
  //     projectId: this.projectId,
  //     memberId: memberId
  //   });
  // }

  getDefaults() {
    this.loading = {
      project: true
    };
    this.project = {
      id: "mock_project_id",
      name: "Mock Project AOSP",
      logoRef: "mocklogo.com",
      description: "description"
    };
    this.sort = {
      alert: {
        arrayProp: "alert",
        sortType: "text",
        sort: null
      },
      company: {
        arrayProp: "org.name",
        sortType: "text",
        sort: null
      },
      product: {
        arrayProp: "product",
        sortType: "text",
        sort: null
      },
      status: {
        arrayProp: "invoices[0].status",
        sortType: "text",
        sort: null
      },
      dues: {
        arrayProp: "annualDues",
        sortType: "number",
        sort: null
      },
      renewal: {
        arrayProp: "renewalDate",
        sortType: "date",
        sort: null
      }
    };
  }

  // sortMembers(prop) {
  //   this.sortService.toggleSort(this.sort, prop, this.project.members);
  // }
}
