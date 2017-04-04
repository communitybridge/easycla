import { Component, ViewChild } from '@angular/core';

import { NavController } from 'ionic-angular';

import { Chart } from 'chart.js';

@Component({
  selector: 'projects-list',
  templateUrl: 'projects-list.html'
})
export class ProjectsListPage {
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


  constructor(public navCtrl: NavController) {
    this.contracts = {
      base: [2, 5, 3],
      additional: [0, 0, 0],
      alert: [0, 0, 0],
    };
    this.invoices = {
      base: [2, 4, 3],
      additional: [1, 0, 0],
      alert: [0, 2, 0],
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
        icon: "https://dummyimage.com/600x250/ffffff/000.png&text=zephyr+logo",
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
        icon: "https://dummyimage.com/600x250/ffffff/000.png&text=todo+logo",
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
        icon: "https://dummyimage.com/600x250/ffffff/000.png&text=openswitch+logo",
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
        icon: "https://dummyimage.com/600x250/ffffff/000.png&text=openchain+logo",
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

  projectSelected(event, project) {
    // That's right, we're pushing to ourselves!
    this.navCtrl.push(ProjectsListPage, {
      project: project
    });
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
        labels: ["Pending", "Member", "LF CoSign"],
        // Base
        datasets: [{
          data: this.contracts.base,
          backgroundColor: this.chartColors.base.background,
          hoverBackgroundColor: this.chartColors.base.hoverBackground
        },
        // Additional
        {
          data: this.contracts.additional,
          backgroundColor: this.chartColors.additional.hoverBackground,
          hoverBackgroundColor: this.chartColors.additional.background
        },
        // Alert
        {
          data: this.contracts.alert,
          backgroundColor: this.chartColors.alert.background,
          hoverBackgroundColor: this.chartColors.alert.hoverBackground
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
