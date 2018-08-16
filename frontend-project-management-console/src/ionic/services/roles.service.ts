import { Injectable } from "@angular/core";
import { Observable } from "rxjs/Observable";
import { KeycloakService } from "./keycloak/keycloak.service";
import { AuthService } from "./auth.service";

@Injectable()
export class RolesService {
  public userAuthenticated: boolean;
  public userRoleDefaults: any;
  public userRoles: any;
  private getDataObserver: any;
  public getData: any;
  private rolesFetched: boolean;

  constructor(
    private keycloak: KeycloakService,
    private authService: AuthService
  ) {
    this.rolesFetched = false;
    this.userRoleDefaults = {
      isAuthenticated: this.authService.isAuthenticated(),
      isPmcUser: true,
      isStaffInc: true,
      isDirectorInc: false,
      isStaffDirect: false,
      isDirectorDirect: false,
      isExec: false,
      isAdmin: true
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
          if (tokenParsed && tokenParsed.roles) {
            this.userRoles = {
              isAuthenticated: this.authService.isAuthenticated(),
              isPmcUser: this.isInArray(tokenParsed.roles, "PMC_LOGIN"),
              isStaffInc: this.isInArray(tokenParsed.roles, "STAFF_STAFF_INC"),
              isDirectorInc: this.isInArray(
                tokenParsed.roles,
                "STAFF_DIRECTOR_INC"
              ),
              isStaffDirect: this.isInArray(
                tokenParsed.roles,
                "STAFF_STAFF_DIRECT"
              ),
              isDirectorDirect: this.isInArray(
                tokenParsed.roles,
                "STAFF_DIRECTOR_DIRECT"
              ),
              isExec: this.isInArray(tokenParsed.roles, "STAFF_EXEC"),
              isAdmin: this.isInArray(tokenParsed.roles, "STAFF_SUPER_ADMIN")
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
