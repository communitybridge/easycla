// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from "@angular/core";
import { Observable } from "rxjs/Observable";
import { AuthService } from "./auth.service";

@Injectable()
export class RolesService {
  public userAuthenticated: boolean;
  public userRoleDefaults: any;
  public userRoles: any;
  private getDataObserver: any;
  public getData: any;
  private rolesFetched: boolean;

  private LF_USERNAME_CLAIM = "https://sso.linuxfoundation.org/claims/username";
  private CLA_PROJECT_ADMIN = "cla-system-admin";

  constructor(
    private authService: AuthService
  ) {
    this.rolesFetched = false;
    this.userRoleDefaults = {
      isAuthenticated: this.authService.isAuthenticated(),
      isPmcUser: false,
      isStaffInc: false,
      isDirectorInc: false,
      isStaffDirect: false,
      isDirectorDirect: false,
      isExec: false,
      isAdmin: false
    };
    this.userRoles = this.userRoleDefaults;
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * This service should ONLY contain methods for user roles
   **/

  //////////////////////////////////////////////////////////////////////////////
  //////////////////////////////////////////////////////////////////////////////

  getUserRolesPromise() {
    console.log("Get UserRole Promise.");
    if (this.authService.isAuthenticated()) {
      return this.authService
        .getIdToken()
        .then(token => {
          return this.authService.parseIdToken(token);
        })
        .then(tokenParsed => {
          if (tokenParsed && tokenParsed[this.LF_USERNAME_CLAIM]) {
            this.userRoles = {
              isAuthenticated: this.authService.isAuthenticated(),
              isPmcUser: false,
              isStaffInc: false,
              isDirectorInc: false,
              isStaffDirect: false,
              isDirectorDirect: false,
              isExec: false,
              isAdmin: false,
            };

            return this.userRoles;
          }

          return this.userRoleDefaults;
        })
        .catch(error => {
          return Promise.resolve(this.userRoleDefaults);
        });
    } else {
      // not authenticated. can't decode token. just return defaults
      return Promise.resolve(this.userRoleDefaults);
    }
  }

  private isInArray(roles, role) {
    for (let i = 0; i < roles.length; i++) {
      if (roles[i].toLowerCase() === role.toLowerCase()) {
        return true;
      }
    }
    return false;
  }

  //////////////////////////////////////////////////////////////////////////////
}
