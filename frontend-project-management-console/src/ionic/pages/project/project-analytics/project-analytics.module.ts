import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectAnalyticsPage } from './project-analytics';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { ProjectHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';
import { Ng2GoogleChartsModule } from 'ng2-google-charts';
import { RoundProgressModule } from 'angular-svg-round-progressbar';

@NgModule({
  declarations: [
    ProjectAnalyticsPage,
  ],
  imports: [
    Ng2GoogleChartsModule,
    RoundProgressModule,
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectAnalyticsPage)
  ],
  entryComponents: [
    ProjectAnalyticsPage,
  ]
})
export class ProjectAnalyticsPageModule {}
