import { Component, ViewChild } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'
import { Chart } from 'chart.js';
// import { ProjectPage } from '../project/project';

@IonicPage({
  segment: 'projects-list'
})
@Component({
  selector: 'projects-list',
  templateUrl: 'projects-list.html'
})
export class ProjectsListPage {
  allProjects: any;
  // projectId: String;
  pushAddProjectPage;
  contracts: {
    base: Array<number>,
    additional: Array<number>,
    alert: Array<number>,
  };
  invoices: {
    base: Array<number>,
    additional: Array<number>,
    alert: Array<number>,
  };
  contractChartColors : {
    base: {
      background: string,
      hoverBackground: string,
    },
    additional: {
      background: string,
      hoverBackground: string,
    },
    alert: {
      background: string,
      hoverBackground: string,
    },
  };
  chartColors: {
    base: {
      background: string,
      hoverBackground: string,
    },
    additional: {
      background: string,
      hoverBackground: string,
    },
    alert: {
      background: string,
      hoverBackground: string,
    },
  };
  meetings: Array<{
    label: string,
    value: string
  }>;
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


  constructor(public navCtrl: NavController, private cincoService: CincoService) {
    this.pushAddProjectPage = 'AddProjectPage';
    this.getDefaults();
  }

  ngOnInit(){
    this.getAllProjects();
  };

  getAllProjects(){
    this.cincoService.getAllProjects().subscribe(response => {
      this.allProjects = response;
    });
  }

  viewProject(projectId){
    this.navCtrl.push('ProjectPage', {
      projectId: projectId
    });
  }

  projectSelected(event, project) {
    this.navCtrl.push('ProjectPage', {
      project: project
    });
  }

  getDefaults(){

    this.contracts = {
      base: [2, 0],
      additional: [0, 5],
      alert: [0, 0],
    };
    this.invoices = {
      base: [2, 4, 3],
      additional: [1, 0, 0],
      alert: [0, 2, 0],
    };

    this.contractChartColors = {
      base: {
        background: "rgba(163,131,107,1)",
        hoverBackground: "rgba(193,141,107,1)",
      },
      additional: {
        background: "rgba(225,170,128,1)",
        hoverBackground: "rgba(225,190,158,1)",
      },
      alert: {
        background: "rgb(235,119,28)",
        hoverBackground: "rgb(213,107,24)",
      },
    };
    this.chartColors = {
      base: {
        background: "rgba(63,103,126,1)",
        hoverBackground: "rgba(50,90,100,1)",
      },
      additional: {
        background: "rgba(23,83,106,1)",
        hoverBackground: "rgba(20,60,80,1)",
      },
      alert: {
        background: "rgb(235,119,28)",
        hoverBackground: "rgb(213,107,24)",
      },
    };
    this.meetings = [
      {
        label: "TODO",
        value: "3/16/2017, 4:00 - 5:00 pst"
      },
      {
        label: "Open Switch",
        value: "3/17/2017, 4:00 - 5:00 pst"
      },
      {
        label: "Zephyr",
        value: "3/18/2017, 4:00 - 5:00 pst"
      },
    ];

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

  ionViewDidLoad() {
    var barOptions_stacked = {
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
                  fontSize:11
              },
              stacked: true
          }]
      },
      legend:{
          display: false
      },
    };

    this.contractsChart = new Chart(this.contractsCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["NEW", "RENEWAL"],
        // Base
        datasets: [{
          data: this.contracts.base,
          backgroundColor: this.contractChartColors.base.background,
          hoverBackgroundColor: this.contractChartColors.base.hoverBackground
        },
        // Additional
        {
          data: this.contracts.additional,
          backgroundColor: this.contractChartColors.additional.hoverBackground,
          hoverBackgroundColor: this.contractChartColors.additional.background
        },
        // Alert
        {
          data: this.contracts.alert,
          backgroundColor: this.contractChartColors.alert.background,
          hoverBackgroundColor: this.contractChartColors.alert.hoverBackground
        },
      ]
      },

      options: barOptions_stacked,
    });

    this.invoicesChart = new Chart(this.invoicesCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["Upcoming", "Delivered", "Paid"],
        // Base
        datasets: [{
          label: 'Base',
          data: this.invoices.base,
          backgroundColor: this.chartColors.base.background,
          hoverBackgroundColor: this.chartColors.base.hoverBackground
        },
        // Additional
        {
          label: 'Additional',
          data: this.invoices.additional,
          backgroundColor: this.chartColors.additional.hoverBackground,
          hoverBackgroundColor: this.chartColors.additional.background
        },
        // Alert
        {
          label: 'Alert',
          data: this.invoices.alert,
          backgroundColor: this.chartColors.alert.background,
          hoverBackgroundColor: this.chartColors.alert.hoverBackground
        },
      ]
      },

      options: barOptions_stacked,
    });

  }

}
