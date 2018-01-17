import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { KeycloakService } from './keycloak/keycloak.service';

@Injectable()
export class RolesService {

  public userAuthenticated: boolean;
  public userRoleDefaults: any;
  public userRoles: any;
  private getDataObserver: any;
  public getData: any;
  private rolesFetched: boolean;

  constructor (
    private keycloak: KeycloakService,
  ) {
    this.rolesFetched = false;
    this.userRoleDefaults = {
      isAuthenticated: this.keycloak.authenticated(),
      isPmcUser: false,
      isStaffInc: false,
      isDirectorInc: false,
      isStaffDirect: false,
      isDirectorDirect: false,
      isExec: false,
      isAdmin: false,
    };
    this.userRoles = this.userRoleDefaults;
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * This service should ONLY contain methods for user roles
  **/

  //////////////////////////////////////////////////////////////////////////////
  //////////////////////////////////////////////////////////////////////////////

  isInArray(roles, role) {
    for(let i=0; i<roles.length; i++) {
      if (roles[i].toLowerCase() === role.toLowerCase()) {
        return true;
      }
    }
    return false;
  }

  getUserRolesPromise() {
    if (this.keycloak.authenticated()) {
      return this.keycloak.getTokenParsed().then((tokenParsed) => {
        if (tokenParsed && tokenParsed.realm_access && tokenParsed.realm_access.roles) {
          this.userRoles = {
            isAuthenticated: this.keycloak.authenticated(),
            isPmcUser: this.isInArray(tokenParsed.realm_access.roles, 'PMC_LOGIN'),
            isStaffInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_INC'),
            isDirectorInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_INC'),
            isStaffDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_DIRECT'),
            isDirectorDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_DIRECT'),
            isExec: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_EXEC'),
            isAdmin: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_SUPER_ADMIN'),
          };
          return this.userRoles;
        }
        return this.userRoleDefaults;
      });
    } else { // not authenticated. can't decode token. just return defaults
      return Promise.resolve(this.userRoleDefaults);
    }
  }

  //////////////////////////////////////////////////////////////////////////////

}
