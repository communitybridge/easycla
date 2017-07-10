import { Component, ViewChild } from '@angular/core';
import { DomSanitizer} from '@angular/platform-browser';
import { NavController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'
import { Chart } from 'chart.js';

@IonicPage({
  segment: 'projects'
})
@Component({
  selector: 'projects-list',
  templateUrl: 'projects-list.html'
})
export class ProjectsListPage {
  loading: any;
  allProjects: any;
  numberOfContracts: {
    new: number,
    renewal: number,
  }
  user: any;
  numberOfInvoices: {
    fewerThan60Days: number,
    fewerThan30Days: number,
    sent: number,
    late: number,
    paidLast30Days: number,
  }
  projects: Array<{
    icon: string,
    title: string,
    datas: Array<{
      label: string,
      value: string,
      status?: string
    }>,
    meetings: Array<{
      label: string,
      value: string
    }>
  }>;

  @ViewChild('contractsCanvas') contractsCanvas;
  contractsChart: any;

  @ViewChild('invoicesCanvas') invoicesCanvas;
  invoicesChart: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private sanitizer: DomSanitizer,
  ) {
    this.getDefaults();
  }

  ngOnInit(){
    this.getAllProjects();
    this.getCurrentUser();
  }

  getAllProjects(){
    this.cincoService.getAllProjects().subscribe(response => {
        this.allProjects = response;
        this.loading.projects = false;
    });
  }

  getCurrentUser(){
    this.cincoService.getCurrentUser().subscribe(response => {
      this.user = response;
      if (response.calendar !== null) {
        this.user.calendar = this.sanitizer.bypassSecurityTrustResourceUrl(response.calendar);
      }
    });
  }

  viewProject(projectId){
    this.navCtrl.setRoot('ProjectPage', {
      projectId: projectId
    });
  }

  projectSelected(event, project) {
    this.navCtrl.push('ProjectPage', {
      project: project
    });
  }

  ionViewDidLoad() {
    let barOptions = this.getBarOptions();
    this.contractsChart = new Chart(this.contractsCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["NEW", "RENEWAL"],
        datasets: [{
             label: '# of Contracts',
             data: [this.numberOfContracts.new, this.numberOfContracts.renewal],
             backgroundColor: [
                 'rgba(163,131,107,1)',
                 'rgba(225,170,128,1)',
             ]
         }]
      },
      options: barOptions
    });

    this.invoicesChart = new Chart(this.invoicesCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["<60 Days", "<30 Days", "SENT", "LATE", "PAID"],
        datasets: [{
             label: '# of Invoices',
             data: [
               this.numberOfInvoices.fewerThan60Days,
               this.numberOfInvoices.fewerThan30Days,
               this.numberOfInvoices.sent,
               this.numberOfInvoices.late,
               this.numberOfInvoices.paidLast30Days],
             backgroundColor: [
                 'rgba(196,221,140,1)',
                 'rgba(171,206,92,1)',
                 'rgba(136,186,22,1)',
                 'rgba(245,166,35,1)',
                 'rgba(65,117,5,1)',
             ]
         }]
      },
      options: barOptions
    });

  }

  getBarOptions(){
    return {
      layout:{
        padding: 20
      },
      responsive: true,
      tooltips: {
          enabled: true
      },
      hover :{
          animationDuration:0
      },
      scales: {
          xAxes: [{
              ticks: {
                  beginAtZero:true,
                  fontSize:11
              },
              scaleLabel:{
                  display:false
              },
              gridLines: {
                display:false,
              },
              stacked: true
          }],
          yAxes: [{
              gridLines: {
                  display:false,
                  color: "#fff",
                  zeroLineColor: "#fff",
                  zeroLineWidth: 0
              },
              ticks: {
                  beginAtZero: true,
                  fontSize:11
              },
              stacked: true,
              display: false,
              barThickness: 300,
          }]
      },
      legend:{
          display: false
      }
    }
  }

  getDefaults(){
    this.loading = {
      charts: true,
      projects: true,
    }

    this.numberOfContracts = {
      new: 15,
      renewal: 50
    };

    this.numberOfInvoices = {
      fewerThan60Days: 50,
      fewerThan30Days: 50,
      sent: 60,
      late: 15,
      paidLast30Days: 15,
    };

    this.user = {
      userId: "",
      email: "",
      roles: [],
      calendar: null,
    }

    this.projects = [
      {
        icon: "assets/test/zephyr-logo.png",
        title: "Zephyr",
        datas: [
          {
            label: "Upcoming",
            value: "2",
          },
          {
            label: "Contracts",
            value: "1",
          },
          {
            label: "Invoices",
            value: "3 (2)",
            status: "alert"
          },
          {
            label: "Paid",
            value: "2"
          },
        ],
        meetings: [
          {
            label: "Zeph Board Meeting",
            value: "3/18/2017, 4:00 - 5:00 pst"
          },
        ],
      },
      {
        icon: "assets/test/todo-logo.png",
        title: "TODO",
        datas: [
          {
            label: "Upcoming",
            value: "2",
          },
          {
            label: "Contracts",
            value: "1",
          },
          {
            label: "Invoices",
            value: "3 (2)",
            status: "alert"
          },
          {
            label: "Paid",
            value: "2"
          },
        ],
        meetings: [
          {
            label: "Board Meeting",
            value: "3/18/2017, 4:00 - 5:00 pst"
          }
        ],
      },
      {
        icon: "assets/test/openswitch-logo.png",
        title: "OpenSwitch",
        datas: [
          {
            label: "Upcoming",
            value: "2",
          },
          {
            label: "Contracts",
            value: "1",
          },
          {
            label: "Invoices",
            value: "3 (2)",
            status: "alert"
          },
          {
            label: "Paid",
            value: "2"
          },
        ],
        meetings: [
          {
            label: "Board Meeting",
            value: "3/18/2017, 4:00 - 5:00 pst"
          }
        ],
      },
      {
        icon: "assets/test/openchain-logo.png",
        title: "OpenChain",
        datas: [
          {
            label: "Upcoming",
            value: "9",
          },
          {
            label: "Contracts",
            value: "1",
          },
          {
            label: "Invoices",
            value: "3 (2)",
            status: "alert"
          },
          {
            label: "Paid",
            value: "2"
          },
        ],
        meetings: [
          {
            label: "Board Meeting",
            value: "3/18/2017, 4:00 - 5:00 pst"
          }
        ],
      },
    ];
  }

  openAccountSettings() {
    this.navCtrl.setRoot('AccountSettingsPage');
  }

}
