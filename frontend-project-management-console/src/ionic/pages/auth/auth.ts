import { Component } from '@angular/core';
import { NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { RolesService } from "../../services/roles.service";

/**
 * Generated class for the AuthPage page.
 *
 * See https://ionicframework.com/docs/components/#navigation for more info on
 * Ionic pages and navigation.
 */

 @Component({
  selector: 'page-auth',
  templateUrl: 'auth.html',
})
export class AuthPage {

  constructor(public navCtrl: NavController, public authService: AuthService, public rolesService: RolesService) {
  }

  ionViewDidEnter() {
    console.log('ionViewDidEnter ProjectsPage');

    setTimeout(() => {
      this.rolesService.getUserRolesPromise().then((userRoles) => {
        if(this.hasAccess(userRoles)) {
          this.navCtrl.setRoot("AllProjectsPage");
        } else {
          this.navCtrl.setRoot("LoginPage");
        }
      });
    }, 2000);
  }

  private hasAccess(userRoles: any): boolean {
    return userRoles.isAuthenticated && userRoles.isPmcUser;
  }

}
