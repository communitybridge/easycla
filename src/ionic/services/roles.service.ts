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
      authenticated: this.keycloak.authenticated(),
      user: false,
      staffInc: false,
      directorInc: false,
      staffDirect: false,
      directorDirect: false,
      exec: false,
      admin: false,
    };

    this.userRoles = this.userRoleDefaults;
    console.log('rolesService userRoles defaults set:');
    console.log(this.userRoles);
    // this.getDataObserver = null;
    // this.getData = Observable.create(observer => {
    //     this.getDataObserver = observer;
    // });
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
      if (this.keycloak.authenticated()) {
        this.keycloak.getTokenParsed().then((tokenParsed) => {
          if (tokenParsed && tokenParsed.realm_access && tokenParsed.realm_access.roles) {
            console.log('tokenParsed:');
            console.log(tokenParsed);
            this.userRoles = {
              authenticated: this.keycloak.authenticated(),
              user: this.isInArray(tokenParsed.realm_access.roles, 'PMC_LOGIN'),
              staffInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_INC'),
              directorInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_INC'),
              staffDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_DIRECT'),
              directorDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_DIRECT'),
              exec: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_EXEC'),
              admin: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_SUPER_ADMIN'),
            };
            this.rolesFetched = true;
            this.getDataObserver.next(this.userRoles);
          }
        });
      }
      // else {
      //   console.log('passing back userRoleDefaults');
      //   this.getDataObserver.next(this.userRoleDefaults); // pass back defaults
      // }

    }
  }

  getUserRolesPromise() {
    if (this.keycloak.authenticated()) {
      return this.keycloak.getTokenParsed().then((tokenParsed) => {
        if (tokenParsed && tokenParsed.realm_access && tokenParsed.realm_access.roles) {
          console.log('getUserRolesPromise tokenParsed:');
          console.log(tokenParsed);
          this.userRoles = {
            authenticated: this.keycloak.authenticated(),
            user: this.isInArray(tokenParsed.realm_access.roles, 'PMC_LOGIN'),
            staffInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_INC'),
            directorInc: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_INC'),
            staffDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_STAFF_DIRECT'),
            directorDirect: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_DIRECTOR_DIRECT'),
            exec: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_EXEC'),
            admin: this.isInArray(tokenParsed.realm_access.roles, 'STAFF_SUPER_ADMIN'),
          };
          return this.userRoles;
        }
        return this.userRoleDefaults;
      });
    } else {
      return Promise.resolve({});
    }
  }

  //////////////////////////////////////////////////////////////////////////////

}
