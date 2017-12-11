import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { AnalyticsService } from '../../../services/analytics.service';
import { DomSanitizer} from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';
import { HostListener } from '@angular/core'

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
@IonicPage({
  segment: 'project/:projectId/analytics'
})
@Component({
  selector: 'project-analytics',
  templateUrl: 'project-analytics.html',
  providers: [CincoService, AnalyticsService]
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
    private analyticsService: AnalyticsService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');

  }

  public columnChartData:any =  {
    chartType: 'ColumnChart',
    dataTable: [
      ['Date', 'Commits'],
      ['', 0]
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
      // colors: ['#7f97b2'],
      colors: ['#fff'],
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
      colors: ['#004421','#429b21'],
      backgroundColor: '#4e92df',
      legend: 'none',
      // is3D: true
    }
  };

  public area2ChartData:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'Page Views'],
      ['', 0]
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

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.getCommitActivity();
    this.getWebsiteDuration();
    this.redrawCharts();
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

  getCommitActivity() {
    let metricType= 'code.commits';
    let groupBy= 'day';
    let tsFrom= '1510430520000';
    let tsTo=   '1514764800000';
    this.analyticsService.getMetrics(metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.columnChartData.dataTable.push([key, value]);
              this.columnChartData = Object.create(this.columnChartData);
            }
          }
        );
      }
    });
  }

  getWebsiteDuration() {
    let metricType= 'website.duration';
    let groupBy= 'week';
    let tsFrom= '1510430520000';
    let tsTo=   '1514764800000';
    this.analyticsService.getMetrics(metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.area2ChartData.dataTable.push([key, value]);
              this.area2ChartData = Object.create(this.area2ChartData);
            }
          }
        );
      }
    });
  }

  redrawCharts() {
    this.columnChartData = Object.create(this.columnChartData);
    this.area2ChartData = Object.create(this.area2ChartData);
    this.pieChartData = Object.create(this.pieChartData);
    this.areaChartData = Object.create(this.areaChartData);
    this.gaugeChartData = Object.create(this.gaugeChartData);
  }

  @HostListener('window:resize', ['$event'])
  onResize(event) {
    event.target.innerWidth;
    this.redrawCharts();
  }

}
