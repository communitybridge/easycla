export function Restricted(restrictions: any) {

  return function (target: Function) {
    target.prototype.ionViewCanEnter = function () {
      console.log(restrictions.roles);
      if (restrictions.roles) {
        if (!this.rolesService) {
          console.warn('[WARNING] this.rolesService is not defined for ' + target.prototype.constructor.name);
          return true; 
        }
        return this.rolesService.getUserRolesPromise().then((userRoles) => {
          let access = true;

          console.log(restrictions.roles);
        
          if (access) {
            return true;
          } else {
            console.log('no access');
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
