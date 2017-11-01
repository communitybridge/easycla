import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectFilesPage } from './project-files';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { ProjectHeaderComponentModule } from '../../../components/project-header/project-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';

@NgModule({
  declarations: [
    ProjectFilesPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectFilesPage)
  ],
  entryComponents: [
    ProjectFilesPage,
  ]
})
export class ProjectFilesPageModule {}
