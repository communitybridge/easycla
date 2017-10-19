import { Component, ViewChild } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';

import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'all-projects-logos'
})
@Component({
  selector: 'all-projects-logos',
  templateUrl: 'all-projects-logos.html'
})
export class AllProjectsLogosPage {

  allProjects: any;
  allProjectsLogos: any[] = [];

  project = new ProjectModel();

  allProjectsWithMembers: any[] = [];

  allProjectsWithLogos: any[] = [];
  projectLogos: any;

  membersCount: any;

  constructor(public navCtrl: NavController, private cincoService: CincoService, private keycloak: KeycloakService) {
    this.getDefaults();
  }

  ionViewCanEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.setRoot('LoginPage');
      this.navCtrl.popToRoot();
    }
    return this.keycloak.authenticated();
  }

  ionViewWillEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.push('LoginPage');
    }
  }

  async ngOnInit(){
    this.getAllProjects();
    console.log("this.allProjectsWithLogos");
    console.log(this.allProjectsWithLogos);
  }

  getAllProjects() {
    this.cincoService.getAllProjects().subscribe(response => {
      if(response) {
        this.allProjects = response;
        for(let eachProject of this.allProjects) {
          this.getAllProjectsLogos(eachProject);
        }
      }
    });
  }

  getAllProjectsLogos(project) {
    this.cincoService.getProjectLogos(project.projectId).subscribe(response => {

      if(response) {
        this.projectLogos = response;
        // Temporary Fix until CINCO returns filename in GET logos response
        for(let eachLogo of this.projectLogos) {
          eachLogo.key = eachLogo.key.split("/");
          eachLogo.key = eachLogo.key[3];
        }
        this.allProjectsWithLogos.push(project, this.projectLogos);
        // for(let eachProject of this.allProjects) {
        //   // this.getProject(eachProject.id);
        // }
      }
      else{
        this.allProjectsWithLogos.push(project);
      }

    });
  }

  // getProject(projectId) {
  //   let getMembers = true;
  //
  //   this.cincoService.getProject(projectId, getMembers).subscribe(response => {
  //     if(response) {
  //       this.allProjectsWithMembers.push(response);
  //     }
  //   });
  //
  //   this.cincoService.getProjectLogos(projectId).subscribe(response => {
  //     // CINCO sample response
  //     // classifier: "main"
  //     // key:"logos/project/a090n000000uQhdAAE/main.png"
  //     // publicUrl:"http://docker.for.mac.localhost:50563/public-media.platform.linuxfoundation.org/logos/project/a090n000000uQhdAAE/main.png
  //     if(response) {
  //
  //       this.projectLogos = response;
  //       // Temporary Fix until CINCO returns filename in GET logos response
  //       for(let eachLogo of this.projectLogos) {
  //         eachLogo.key = eachLogo.key.split("/");
  //         eachLogo.key = eachLogo.key[3];
  //       }
  //
  //       this.allProjectsWithLogos.push(this.projectLogos);
  //       console.log(this.allProjectsWithLogos);
  //
  //     }
  //
  //   });
  //
  // }

  ionViewDidLoad() {

  }

  getDefaults() {

  }

}
