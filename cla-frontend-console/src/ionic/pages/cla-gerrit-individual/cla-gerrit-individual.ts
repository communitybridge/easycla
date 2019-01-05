import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { RolesService } from "../../services/roles.service";
import { AuthService } from "../../services/auth.service";
import { Restricted } from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: 'cla/gerrit/project/:projectId/individual'
})
@Component({
  selector: 'cla-gerrit-individual',
  templateUrl: 'cla-gerrit-individual.html'
})
export class ClaGerritIndividualPage {
  projectId: string;
  project: any;
  user: any;
  signatureIntent: any;
  activeSignatures: boolean = true; // we assume true until otherwise
  signature: any;

  userRoles: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService,
    private rolesService: RolesService,
    private authService: AuthService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    localStorage.setItem("gerritProjectId", this.projectId);
    localStorage.setItem("gerritClaType", "ICLA");
  }

  getDefaults() {
    this.userRoles = this.rolesService.userRoleDefaults;

    this.project = {
      project_name: "",
    };
    this.signature = {
      sign_url: "",
    };
  }

  ngOnInit() {
    this.getProject(this.projectId);
    this.getUserInfo();
  }

  ionViewCanEnter(){
    if(!this.authService.isAuthenticated){
      setTimeout(()=>this.navCtrl.setRoot('LoginPage'))
    }
    return this.authService.isAuthenticated
  }


  ngAfterViewInit() {
  }

  getUserInfo() {
    this.authService.getUserInfo().then(res => {
      this.user = res;
      console.log(this.user.email);
      console.log(this.user);
    })
  }
  
  getOrCreateUser() {
    this.claService.getUserByEmail(this.userEmail).subscribe(user => {
      if(user.errors != null) {
        //create user since user doesnt exist in the db. 
        this.claService.postUser( 
          { user_email: this.user.email, 
            user_name: this.user.name } );

      }
    })
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
      if (!this.project.logoRef) {
        this.project.logoRef = "https://dummyimage.com/200x100/bbb/fff.png&text=+";
      }
    });
  }
  
  openClaAgreement() {
    if (!this.signature.sign_url) {
      return;
    }
    window.open(this.signature.sign_url, '_blank');
  }
}
