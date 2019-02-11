// import { Component, ViewChild } from '@angular/core';
import { Component } from "@angular/core";
// import { DomSanitizer} from '@angular/platform-browser';
import { NavController, IonicPage } from "ionic-angular";
// import { CincoService } from "../../services/cinco.service";
// import { Chart } from 'chart.js';
import { ClaService } from "../../services/cla.service";
import { FilterService } from "../../services/filter.service";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated", "isPmcUser"]
})
@IonicPage({
  name: "AllProjectsPage",
  segment: "projects"
})
@Component({
  selector: "all-projects",
  templateUrl: "all-projects.html"
})
export class AllProjectsPage {
  loading: any;
  expand: any;
  projectSectors: any;
  // userHasCalendar: boolean;
  allProjects: any;
  allFilteredProjects: any;
  // numberOfContracts: {
  //   new: number,
  //   renewal: number,
  // }
  // user: any;
  // numberOfInvoices: {
  //   fewerThan60Days: number,
  //   fewerThan30Days: number,
  //   sent: number,
  //   late: number,
  //   paidLast30Days: number,
  // }
  // projects: Array<{
  //   icon: string;
  //   title: string;
  //   datas: Array<{
  //     label: string;
  //     value: string;
  //     status?: string;
  //   }>;
  //   meetings: Array<{
  //     label: string;
  //     value: string;
  //   }>;
  // }>;

  // industryFilterValues: Array<{
  //   key: string;
  //   prettyValue: string;
  // }> = [];

  managersFilterValues: any = [];

  userRoles: any;

  // @ViewChild('contractsCanvas') contractsCanvas;
  // contractsChart: any;
  //
  // @ViewChild('invoicesCanvas') invoicesCanvas;
  // invoicesChart: any;

  constructor(
    public navCtrl: NavController,
    // private cincoService: CincoService,
    // private sanitizer: DomSanitizer,
    private claService: ClaService,
    private rolesService: RolesService,
    private filterService: FilterService
  ) {
    this.getDefaults();
  }

  async ngOnInit() {
    // this.getIndustries();
    // this.getCurrentUser();
    this.getAllProjectFromSFDC();
  }

  getAllProjectFromSFDC() {
    var self = this;
    this.claService.getAllProjectsFromSFDC({
        "400":function(rawErrorObject) {
          /* 
            realize you get the whole error object here you can further 
            breakdown possible paths depending on the actual output
          */
          // console.log(rawErrorObject);
          self.navCtrl.setRoot("LoginPage");
        },
        "401":function(rawErrorObject) {
          self.navCtrl.setRoot("LoginPage");
        }
    }).subscribe(response => {
      this.allProjects = response;
      this.allFilteredProjects = this.filterService.resetFilter(
        this.allProjects
      );
      this.loading.projects = false;
    });
    // we may need to do some filter later instead of showing all projects
  }

  // getAllProjects() {
  //   this.cincoService.getAllMockProjects().subscribe(response => {
  //     this.allProjects = response;
  //     for (let eachProject of this.allProjects) {
  //       // After uploading a logo, Cinco will provide same name,
  //       // so a refresh to the image needs to be forced.
  //       // This is to refresh an image that have same URL
  //       if (eachProject.config.logoRef) {
  //         eachProject.config.logoRef += "?" + new Date().getTime();
  //       }
  //       // Currently PMs are returned with their KC IDs
  //       // This translates to LF IDs
  //       if (eachProject.config.programManagers.length > 0) {
  //         for (let eachManager of eachProject.config.programManagers) {
  //           this.cincoService.getUser(eachManager).subscribe(response => {
  //             if (response) {
  //               eachProject.managers.push(response.lfId);
  //               // Prevent dupes in PMs filter list
  //               if (
  //                 !this.managersFilterValues.some(lfId => lfId == response.lfId)
  //               ) {
  //                 this.managersFilterValues.push(response.lfId);
  //               }
  //               this.allFilteredProjects = this.filterService.resetFilter(
  //                 this.allProjects
  //               );
  //             }
  //           });
  //         }
  //       }
  //     }
  //     this.allFilteredProjects = this.filterService.resetFilter(
  //       this.allProjects
  //     );
  //     this.loading.projects = false;
  //   });
  // }

  // getCurrentUser() {
  //   this.cincoService.getCurrentUser().subscribe(response => {
  //     this.user = response;
  //     if (response.hasOwnProperty('calendar') && response.calendar) {
  //       this.userHasCalendar = true;
  //       this.user.calendar = this.sanitizer.bypassSecurityTrustResourceUrl(response.calendar);
  //     }
  //   });
  // }

  viewProject(projectId) {
    this.navCtrl.setRoot("ProjectPage", {
      projectId: projectId
    });
  }

  viewProjectCLA(projectId) {
    this.navCtrl.setRoot("ProjectClaPage", {
      projectId: projectId
    });
  }

  // projectSelected(event, project) {
  //   this.navCtrl.push('ProjectPage', {
  //     project: project
  //   });
  // }

  // ionViewDidLoad() {
  //   let barOptions = this.getBarOptions();
  //   this.contractsChart = new Chart(this.contractsCanvas.nativeElement, {
  //     type: 'bar',
  //     data: {
  //       labels: ["NEW", "RENEWAL"],
  //       datasets: [{
  //            label: '# of Contracts',
  //            data: [this.numberOfContracts.new, this.numberOfContracts.renewal],
  //            backgroundColor: [
  //                'rgba(163,131,107,1)',
  //                'rgba(225,170,128,1)',
  //            ]
  //        }]
  //     },
  //     options: barOptions
  //   });
  //
  //   this.invoicesChart = new Chart(this.invoicesCanvas.nativeElement, {
  //     type: 'bar',
  //     data: {
  //       labels: ["<60 Days", "<30 Days", "SENT", "LATE", "PAID"],
  //       datasets: [{
  //            label: '# of Invoices',
  //            data: [
  //              this.numberOfInvoices.fewerThan60Days,
  //              this.numberOfInvoices.fewerThan30Days,
  //              this.numberOfInvoices.sent,
  //              this.numberOfInvoices.late,
  //              this.numberOfInvoices.paidLast30Days],
  //            backgroundColor: [
  //                'rgba(196,221,140,1)',
  //                'rgba(171,206,92,1)',
  //                'rgba(136,186,22,1)',
  //                'rgba(245,166,35,1)',
  //                'rgba(65,117,5,1)',
  //            ]
  //        }]
  //     },
  //     options: barOptions
  //   });
  //
  // }

  // getBarOptions(){
  //   return {
  //     layout:{
  //       padding: 20
  //     },
  //     responsive: true,
  //     tooltips: {
  //         enabled: true
  //     },
  //     hover :{
  //         animationDuration:0
  //     },
  //     scales: {
  //         xAxes: [{
  //             ticks: {
  //                 beginAtZero:true,
  //                 fontSize:11
  //             },
  //             scaleLabel:{
  //                 display:false
  //             },
  //             gridLines: {
  //               display:false,
  //             },
  //             stacked: true
  //         }],
  //         yAxes: [{
  //             gridLines: {
  //                 display:false,
  //                 color: "#fff",
  //                 zeroLineColor: "#fff",
  //                 zeroLineWidth: 0
  //             },
  //             ticks: {
  //                 beginAtZero: true,
  //                 fontSize:11
  //             },
  //             stacked: true,
  //             display: false,
  //             barThickness: 300,
  //         }]
  //     },
  //     legend:{
  //         display: false
  //     }
  //   }
  // }

  getDefaults() {
    this.userRoles = this.rolesService.userRoleDefaults;

    // this.userHasCalendar = false;

    this.loading = {
      // charts: true,
      projects: true
    };

    this.expand = {};

    // this.numberOfContracts = {
    //   new: 15,
    //   renewal: 50
    // };
    //
    // this.numberOfInvoices = {
    //   fewerThan60Days: 50,
    //   fewerThan30Days: 50,
    //   sent: 60,
    //   late: 15,
    //   paidLast30Days: 15,
    // };

    // this.user = {
    //   userId: "",
    //   email: "",
    //   roles: [],
    //   calendar: null,
    // }
  }

  // openAccountSettings() {
  //   this.navCtrl.setRoot('AccountSettingsPage');
  // }

  // getIndustries() {
  //   this.cincoService.getProjectSectors().subscribe(response => {
  //     this.projectSectors = response;
  //     for (let key in this.projectSectors) {
  //       if (this.projectSectors.hasOwnProperty(key)) {
  //         let industry = {
  //           key: key,
  //           prettyValue: this.projectSectors[key]
  //         };
  //         this.industryFilterValues.push(industry);
  //       }
  //     }
  //   });
  // }

  filterAllProjects(projectProperty, keyword) {
    if (keyword == "NO_FILTER") {
      this.allFilteredProjects = this.filterService.resetFilter(
        this.allProjects
      );
    } else {
      this.allFilteredProjects = this.filterService.filterAllProjects(
        this.allProjects,
        projectProperty,
        keyword
      );
    }
  }

  toggleExpand(index) {
    this.expand[index] = !this.expand[index];
  }
}
