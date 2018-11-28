import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectGroupsCreatePage } from './project-groups-create';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { ProjectHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';

@NgModule({
  declarations: [
    ProjectGroupsCreatePage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectGroupsCreatePage)
  ],
  entryComponents: [
    ProjectGroupsCreatePage,
  ]
})
export class ProjectGroupsCreatePageModule {}
