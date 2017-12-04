import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { KeycloakService } from './keycloak/keycloak.service';

@Injectable()
export class RolesService {

  public userAuthenticated: boolean;
  public userRoleDefaults: any;
  private userRoles: any;
  private getDataObserver: any;
  public getData: any;
  private rolesFetched: boolean;

  constructor (
    private keycloak: KeycloakService,
  ) {
    this.rolesFetched = false;
    this.userRoleDefaults = {
      isAuthenticated: this.keycloak.authenticated(),
      isUser: false,
      isProgramManager: false,
      isProgramManagerAdmin: false,
      isAdmin: false,
      isStaffInc: false,
    };
    this.getDataObserver = null;
    this.getData = Observable.create(observer => {
        this.getDataObserver = observer;
    });
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

  getUserRoles() {
    if (this.rolesFetched) {
      this.getDataObserver.next(this.userRoles);
    } else {
      this.keycloak.getTokenParsed().then((tokenParsed) => {
        if (tokenParsed) {
          let isUser = this.isInArray(tokenParsed.realm_access.roles, 'PMC_LOGIN');
          let isProgramManager = this.isInArray(tokenParsed.realm_access.roles, 'PROGRAM_MANAGER');
          let isProgramManagerAdmin = this.isInArray(tokenParsed.realm_access.roles, 'PMC_PROGRAM_MANAGER_ADMIN');
          let isAdmin = this.isInArray(tokenParsed.realm_access.roles, 'STAFF_SUPER_ADMIN');
          let isStaffInc = this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_INC');
          this.userRoles = {
            isAuthenticated: this.keycloak.authenticated(),
            isUser: isUser,
            isProgramManager: isProgramManager,
            isProgramManagerAdmin: isProgramManagerAdmin,
            isAdmin: isAdmin,
            isStaffInc: isStaffInc
          };
          this.rolesFetched = true;
          this.getDataObserver.next(this.userRoles);
        }
      });
    }
  }

  //////////////////////////////////////////////////////////////////////////////

}
