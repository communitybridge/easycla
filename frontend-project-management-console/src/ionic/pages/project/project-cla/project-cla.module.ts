import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectClaPage } from './project-cla';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { ProjectHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';

@NgModule({
  declarations: [
    ProjectClaPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectClaPage)
  ],
  entryComponents: [
    ProjectClaPage,
  ]
})
export class ProjectClaPageModule {}
