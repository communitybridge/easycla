import { Component } from "@angular/core";
import { NavController, IonicPage, NavParams } from "ionic-angular";
import { KeycloakService } from "../../services/keycloak/keycloak.service";
import { RolesService } from "../../services/roles.service";
import { AuthService } from "../../services/auth.service";

@IonicPage({
  name: "LoginPage",
  segment: "login/:return"
})
@Component({
  selector: "login",
  templateUrl: "login.html"
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
    public authService: AuthService
  ) {
    let dataString = this.navParams.get("return");
    try {
      this.data = JSON.parse(dataString);
      console.log(this.data);
      this.returnData = true;
    } catch (e) {
      this.returnData = false;
    }
    this.userRoles = this.rolesService.userRoles;
    this.rolesService.getUserRolesPromise().then(userRoles => {
      this.userRoles = userRoles;
      this.checkPageReturn();
    });
  }

  checkPageReturn() {
    this.canAccess = this.hasAccess();
    if (this.canAccess && this.returnData) {
      if (this.data.page) {
        console.log("can be accessed to page:");
        if (this.data.params) {
          console.log("page: " + this.data.page);
          console.log("page params: " + this.data.params);
          this.navCtrl.setRoot(this.data.page, this.data.params);
        } else {
          console.log(this.data.page);
          this.navCtrl.setRoot(this.data.page);
        }
      }
    }
  }

  hasAccess() {
    console.log("check if hasAccess");
    console.log("data.roles: " + this.data.roles);
    if (this.data && this.data.roles) {
      for (let role of this.data.roles) {
        if (!this.userRoles[role]) {
          return false;
        }
      }
    }
    return true;
  }

  login() {
    this.authService.login();
  }

  logout() {
    this.authService.logout();
  }
}
