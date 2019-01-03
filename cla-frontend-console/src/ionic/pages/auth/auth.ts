import { Component } from '@angular/core';
import { NavController } from 'ionic-angular';
import { AuthService } from "../../services/auth.service";
import { RolesService } from "../../services/roles.service";

 @Component({
  selector: 'page-auth',
  templateUrl: 'auth.html',
})
export class AuthPage {

  userRoles: any;

  constructor(public navCtrl: NavController, public authService: AuthService, public rolesService: RolesService) {
  }

  ionViewDidEnter() {
    console.log('ionViewDidEnter AuthPage');

    setTimeout(() => {
      this.rolesService.getUserRolesPromise().then((userRoles) => {
      });
    }, 2000); 
  }

  private hasAccess(userRoles: any): boolean {
    return userRoles.isAuthenticated && userRoles.isPmcUser;
  }

}
