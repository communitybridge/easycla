import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { DomSanitizer} from '@angular/platform-browser';

@IonicPage({
  segment: 'project/:projectId/analytics'
})
@Component({
  selector: 'project-analytics',
  templateUrl: 'project-analytics.html',
  providers: [CincoService]
})
export class ProjectAnalyticsPage {

  projectId: string;

  hasAnalyticsUrl: boolean;
  analyticsUrl: any;
  sanitizedAnalyticsUrl: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController
  ) {
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  public columnChartData:any =  {
    chartType: 'ColumnChart',
    dataTable: [
      ['Date', 'Commits'],
      ['11/10', 700],
      ['11/11', 300],
      ['11/12', 400],
      ['11/13', 500],
      ['11/14', 600],
      ['11/15', 800]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of commits', minValue: 0, maxValue: 15},
      colors: ['#7f97b2'],
      backgroundColor: '#4e92df',
      legend: 'none',
      // is3D: true
    }
  };


  public areaChartData:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'PR1', 'PR2'],
      ['11/10', 70, 30],
      ['11/11', 30, 40],
      ['11/12', 40, 70],
      ['11/13', 50, 20],
      ['11/14', 60, 50],
      ['11/15', 80, 10]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of PRs', minValue: 0, maxValue: 15},
      colors: ['#9344dd', '#abab45'],
      backgroundColor: '#4e92df',
      legend: 'none',
      // is3D: true
    }
  };


  public pieChartData:any =  {
    chartType: 'PieChart',
    dataTable: [
      ['Date', 'Commits'],
      ['11/10', 82],
      ['11/11', 18]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of PRs', minValue: 0, maxValue: 15},
      colors: ['#9b9b9b'],
      backgroundColor: '#4e92df',
      legend: 'none',
      // is3D: true
    }
  };

  public area2ChartData:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'Page Views'],
      ['11/10', 20],
      ['11/11', 60],
      ['11/12', 70],
      ['11/13', 20],
      ['11/14', 40],
      ['11/15', 90]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: 'Page Views (in thousands)', minValue: 0, maxValue: 15},
      colors: ['#95c2e2'],
      backgroundColor: '#4e92df',
      legend: 'none',
      // is3D: true
    }
  };


  public tableChartData:any =  {
    chartType: 'Table',
    dataTable: [
      ['Name', 'Email Address', 'PR submit date'],
      ['Nick Young', 'nick@dubs.com', 11/29/17],
      ['Patrick MacCaw', 'patrick@dubs.com', 11/29/17],
      ['David West', 'david@dubs.com', 11/28/17],
      ['Javale MacGee', 'javale@dubs.com', 11/27/17],
      ['Shaun Livingston', 'shaun@dubs.com', 11/26/17],
      ['Andre Iguodala', 'andre@dubs.com', 11/25/17]
    ],
    formatters: [
      {
        columns: [1, 2],
        type: 'NumberFormat',
        options: {
          prefix: '&euro;', negativeColor: 'red', negativeParens: true
        }
      }
    ],
    options: {title: 'Countries', allowHtml: true}
  };


  // public gaugeChartData:any =  {
  //   chartType: 'Gauge',
  //   dataTable: [
  //     ['Date', 'PR'],
  //     ['11/10', 32],
  //   ],
  //   options: {
  //     hAxis: {
  //       textStyle:{ color: '#ffffff'},
  //       gridlines: {
  //         color: "#FFFFFF"
  //       },
  //       baselineColor: '#FFFFFF'
  //     },
  //     vAxis: {title: 'Page Views (in thousands)', minValue: 0, maxValue: 15},
  //     colors: ['#95c2e2'],
  //     backgroundColor: '#4e92df',
  //     legend: 'none',
  //     // is3D: true
  //   }
  // };


  public gaugeChartData:any =  {
    chartType: 'Gauge',
    dataTable: [
      ['Label', 'Value'],
      ['PRs', 32]
    ],
    options: {
      animation: {easing: 'out'},
      // width: 150, height: 150,
      greenFrom: 0, greenTo: 32,
      minorTicks: 1,
      min: 0, max: 100,
      majorTicks: ['0', '20', '40', '60', '80', '100'],
      greenColor: '#d0e9c6'
    }
  };


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

  ngOnInit() {
    this.getProjectConfig(this.projectId);
  }

  getDefaults() {

  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        let projectConfig = response;
        if(projectConfig.analyticsUrl) {
          this.analyticsUrl = projectConfig.analyticsUrl;
          this.sanitizedAnalyticsUrl = this.domSanitizer.bypassSecurityTrustResourceUrl(this.analyticsUrl);
          this.hasAnalyticsUrl = true;
        }
        else{
          this.hasAnalyticsUrl = true;
        }
      }
    });
  }

  openAnaylticsConfigModal(projectId) {
    let modal = this.modalCtrl.create('AnalyticsConfigModal', {
      projectId: projectId,
    });
    modal.onDidDismiss(analyticsUrl => {
      if(analyticsUrl){
        this.analyticsUrl = analyticsUrl;
        this.hasAnalyticsUrl = true;
        this.sanitizedAnalyticsUrl = this.domSanitizer.bypassSecurityTrustResourceUrl(this.analyticsUrl);
      }
    });
    modal.present();
  }

}
