// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from "@angular/core";
import { NavController, IonicPage } from "ionic-angular";
import { ClaService } from "../../services/cla.service";
import { FilterService } from "../../services/filter.service";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated", "isPmcUser"]
})
@IonicPage({
  name: "AllProjectsPage",
  segment: "projects"
})
@Component({
  selector: "all-projects",
  templateUrl: "all-projects.html"
})
export class AllProjectsPage {
  loading: any;
  projectSectors: any;
  allProjects: any;
  allFilteredProjects: any;
  userRoles: any;
  errorMessage = null;
  errorStatus = null;

  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    private rolesService: RolesService,
    private filterService: FilterService
  ) {
    this.getDefaults();
  }

  async ngOnInit() {
    this.getAllProjectFromSFDC();
  }

  getAllProjectFromSFDC() {
    this.claService.getAllProjectsFromSFDC().subscribe(response => {
      this.allProjects = this.sortProjects(response);
      this.allFilteredProjects = this.filterService.resetFilter(
        this.allProjects
      );
      this.loading.projects = false;
    }, (error) => this.handleErrors(error));
  }

  sortProjects(projects) {
    if (projects == null || projects.length == 0) {
      return projects;
    }

    return projects.sort((a, b) => {
      return a.name.localeCompare(b.name);
    });
  }

  handleErrors (error) {
    this.setLoadingSpinner(false);
    this.errorStatus = error.status;

    switch (error.status) {
      case 401:
        this.errorMessage = `Your session may have expired or you don't have permissions to see any projects.`;
        break;

      default:
        this.errorMessage = `An unknown error has occurred when retrieving the projects`;
    }
  }

  viewProjectCLA(projectId) {
    this.navCtrl.setRoot("ProjectClaPage", {
      projectId: projectId
    });
  }

  getDefaults() {
    this.userRoles = this.rolesService.userRoleDefaults;

    this.setLoadingSpinner(true);
  }

  setLoadingSpinner (value) {
    this.loading = {
      projects: value
    };
  }

  filterAllProjects(projectProperty, keyword) {
    if (keyword == "NO_FILTER") {
      this.allFilteredProjects = this.filterService.resetFilter(
        this.allProjects
      );
    } else {
      this.allFilteredProjects = this.filterService.filterAllProjects(
        this.allProjects,
        projectProperty,
        keyword
      );
    }
  }

  /**
   * Opens the access page in a new window
   */
  openAccessPage() {
    window.open('https://app.gitbook.com/@lf-docs-linux-foundation/s/easycla/getting-started/get-access-to-easycla', "_blank");
  }
}
