import { Component } from '@angular/core';
import { NavController, IonicPage, NavParams } from 'ionic-angular';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { RolesService } from '../../services/roles.service';

@IonicPage({
  name: 'LoginPage',
  segment: 'login/:return'
})
@Component({
  selector: 'login',
  templateUrl: 'login.html'
})
export class LoginPage {

  returnData: boolean;
  data: any;
  userRoles: any;
  canAccess: boolean;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    // console.log('login page loaded');
    let dataString = this.navParams.get('return');
    // console.log('data');
    // console.log(dataString);
    try {
      this.data = JSON.parse(dataString);
      this.returnData = true;
    } catch (e) {
      this.returnData = false;
    }
    this.userRoles = this.rolesService.userRoles;
    console.log('login page userroles');
    console.log(this.userRoles);
    this.checkPageReturn();
    // this.rolesService.getData.subscribe((userRoles) => {
    //   this.userRoles = userRoles;
    //   this.checkPageReturn();
    // });
    // this.rolesService.getUserRoles();
  }

  checkPageReturn() {
    console.log('check page return');
    this.canAccess = this.hasAccess();
    console.log(this.canAccess);
    console.log(this.data);
    if (this.canAccess && this.returnData) {
      if (this.data.page) {
        if (this.data.params) {
          this.navCtrl.setRoot(this.data.page, this.data.params);
        } else {
          this.navCtrl.setRoot(this.data.page);
        }
      }
    }
  }

  hasAccess() {
    if (this.data && this.data.roles) {
      console.log('hasAccess: roles required for page:');
      console.log(this.data.roles);
      for (let role of this.data.roles) {
        console.log('restricted role in userRoles:');
        console.log(this.userRoles[role]);
        if (!this.userRoles[role]) {
          console.log('false');
          return false;
        }
      }
    }
    return true;
  }

  login() {
    this.keycloak.login();
  }

  logout() {
    this.keycloak.logout();
  }

}
