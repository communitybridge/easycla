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

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {

    this.getCommitActivity();
    this.getcommitsDistribution();
    this.getIssuesStatus();
    this.getPrsPipeline();
    this.getIssuesActivity();
    this.getPrsActivity();
    this.getPageViews();

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
    let index = 'hyperledger';
    let metricType = 'code.commits';
    let groupBy = 'day';
    let tsFrom = '1510430520000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.commitsActivityChart.dataTable.push([key, value]);
              this.commitsActivityChart = Object.create(this.commitsActivityChart);
            }
          }
        );
      }
    });
  }

  getcommitsDistribution() {
    console.log("TODO");
  }

  getIssuesStatus() {
    console.log("TODO");
  }

  getPrsPipeline() {
    let index = 'hyperledger';
    let metricType = 'prs.open';
    let groupBy = 'day';
    let tsFrom = '1510430520000';
    let tsTo =   '1514764800000';
    let sumOpenPRs = 0;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            sumOpenPRs = sumOpenPRs + value;
          }
        );
        this.prsPipelineChart.dataTable.push(['PRs', sumOpenPRs]);
        this.prsPipelineChart = Object.create(this.prsPipelineChart);
      }
    });
  }

  getIssuesActivity() {
    console.log("TODO");
  }

  getPrsActivity() {
    let index = 'hyperledger';
    let metricType = 'prs.open';
    let groupBy = 'month';
    let tsFrom = '1388534400000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            this.prsActivityChart.dataTable.push([key, value, 0]);
            this.prsActivityChart = Object.create(this.prsActivityChart);
          }
        );
      }
    });
  }

  getPageViews() {
    let index = 'hyperledger';
    let metricType = 'website.duration';
    let groupBy = 'week';
    let tsFrom = '1510430520000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.pageViewsChart.dataTable.push([key, value]);
              this.pageViewsChart = Object.create(this.pageViewsChart);
            }
          }
        );
      }
    });
  }

  redrawCharts() {
    this.commitsActivityChart = Object.create(this.commitsActivityChart);
    this.commitsDistributionChart = Object.create(this.commitsDistributionChart);
    this.issuesStatusChart = Object.create(this.issuesStatusChart);
    this.prsPipelineChart = Object.create(this.prsPipelineChart);
    this.issuesActivityChart = Object.create(this.issuesActivityChart);
    this.prsActivityChart = Object.create(this.prsActivityChart);
    this.pageViewsChart = Object.create(this.pageViewsChart);
  }

  @HostListener('window:resize', ['$event'])
  onResize(event) {
    event.target.innerWidth;
    this.redrawCharts();
  }


  public commitsActivityChart:any =  {
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
      vAxis: {title: '# of commits'},
      colors: ['#fff'],
      backgroundColor: '#4e92df',
      legend: 'none',
    }
  };


  public issuesStatusChart:any =  {
    chartType: 'BarChart',
    dataTable: [
      ['Status', 'Issues'],
      ['News', 1],
      ['Open', 2],
      ['Closed', 3],
      ['Invalid', 4]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of commits'},
      colors: ['#ff88ff'],
      backgroundColor: '#4e92df',
      legend: 'none',
    }
  };


  public prsActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'PRs Open', 'PRs Merged'],
      ['0', 0, 0]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of PRs'},
      colors: ['#9344dd', '#abab45'],
      backgroundColor: '#4e92df',
      legend: 'none'
    }
  };


  public issuesActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'PRs Open', 'PRs Merged'],
      ['11/5', 3, 1],
      ['11/6', 7, 3],
      ['11/7', 3, 4],
      ['11/8', 8, 6],
      ['11/9', 6, 3],
      ['11/10', 2, 1],
      ['11/11', 8, 5]
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of PRs'},
      colors: ['#9344dd', '#abab45'],
      backgroundColor: '#4e92df',
      legend: 'none'
    }
  };


  public commitsDistributionChart:any =  {
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
      legend: 'none'
    }
  };

  public pageViewsChart:any =  {
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
      legend: 'none'
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

  public prsPipelineChart:any =  {
    chartType: 'Gauge',
    dataTable: [
      ['Label', 'Value']
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

}
