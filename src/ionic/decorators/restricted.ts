export function Restricted(restrictions: any) {

  return function (target: Function) {
    console.log('target');
    console.log(target.prototype);
    target.prototype.ionViewCanEnter = function () {
      console.log('restricted roles:');
      console.log(restrictions);
      if (restrictions.roles) {
        console.log('in restrictions roles2');
        console.log(this.rolesService);
        if (!this.rolesService) {
          console.warn('[WARNING] this.rolesService is not defined for ' + target.prototype.constructor.name);
          return true; // Let's not break everything in case we forgot... for now
        }
        return this.rolesService.getUserRolesPromise().then((userRoles) => {
          console.log('restricted roles:');
          console.log(restrictions);
          console.log('promise roles');
          console.log(userRoles);
          let access = true;
          for (let role of restrictions.roles) {
            console.log('restricted role in userRoles:');
            console.log(userRoles[role]);
            if (!userRoles[role]) {
              access = false;
              console.log('false');
              break; // TODO: this doesn't seem to be breaking the for loop
            }
          }

          if (access) {
            return true;
          } else {
            console.log('push to login page with data');
            console.log('current page');
            let navObject = {
              page: target.prototype.constructor.name,
              params: {},
              roles: restrictions.roles,
            };
            console.log(target.prototype.constructor.name);
            if (this.navParams) {
              navObject.params = this.navParams.data;
              console.log('navParams');
              console.log(this.navParams.data);
            }
            console.log('setRoot loginpage2');
            let navString = JSON.stringify(navObject);
            window.location.hash = '#/login/' + navString;
            window.location.reload(true);

            console.log('access denied.')
            return false;
          }
        });
      } else { // no other restrictions implemented yet
        console.log('no roles set. no restrictions');
        return true;
      }
    }

  };

}
