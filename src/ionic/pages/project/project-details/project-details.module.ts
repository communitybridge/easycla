import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectDetailsPage } from './project-details';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { ProjectHeaderComponentModule } from '../../../components/project-header/project-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';

@NgModule({
  declarations: [
    ProjectDetailsPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectDetailsPage)
  ],
  entryComponents: [
    ProjectDetailsPage,
  ]
})
export class ProjectDetailsPageModule {}
