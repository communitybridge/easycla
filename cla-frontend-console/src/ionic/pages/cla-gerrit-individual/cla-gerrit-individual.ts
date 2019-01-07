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
  userId: string;
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
    // retrieve userInfo from auth0 service
    this.authService.getUserInfo().then(res => {
      this.user = res;
      // retrieve existing userId by email, or create one for Gerrit 
      // For users who has an LFID but has never used the CLA app before.
      let data = {
        user_email: this.user.email,
        user_name: this.user.nickname
      }
      this.claService.postOrGetUserForGerrit(data).subscribe(user => {
          this.userId = user.user_id;
          // get signatureIntent object, similar to the Github flow. 
          this.postSignatureRequest();
      })
    })
  }
  
  postSignatureRequest() {
    let signatureRequest = {
      'project_id': this.projectId,
      'user_id': this.userId,
      'return_url_type': "Gerrit",
    };
    this.claService.postIndividualSignatureRequest(signatureRequest).subscribe(response => {
      // returns {
      //   user_id:
      //   signature_id:
      //   project_id:
      //   sign_url: docusign.com/some-docusign-url
      // }
      this.signature = response;
    });
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
