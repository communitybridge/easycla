// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

export function Restricted(restrictions: any) {

  return function (target: Function) {
    target.prototype.ionViewCanEnter = function () {
      if (restrictions.roles) {
        if (!this.rolesService) {
          console.warn('[WARNING] this.rolesService is not defined for ' + target.prototype.constructor.name);
          return true; // Let's not break everything in case we forgot... for now
        }
        return this.rolesService.getUserRolesPromise().then((userRoles) => {
          let access = true;
          for (let role of restrictions.roles) {
            if (!userRoles.hasOwnProperty(role)) {
              console.warn('[WARNING] "' + role + '" is not a defined user role');
            }
            if (!userRoles[role]) {
              access = false;
              break; // TODO: this doesn't seem to be breaking the for loop
            }
          }

          if (access) {
            return true;
          } else {
            console.log('no access');
            // let navObject = {
            //   page: target.prototype.constructor.name,
            //   params: {},
            //   roles: restrictions.roles,
            // };
            // if (this.navParams) {
            //   navObject.params = this.navParams.data;
            // }
            // let navString = JSON.stringify(navObject);
            // window.location.hash = '#/login/';
            window.location.hash = '#/login';
            window.location.reload(true);
            return false;
          }
        });
      } else { // no other restrictions implemented yet
        return true;
      }
    }

  };

}
