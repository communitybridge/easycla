import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectGroupsPage } from './project-groups';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { ProjectHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';

@NgModule({
  declarations: [
    ProjectGroupsPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ProjectHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectGroupsPage)
  ],
  entryComponents: [
    ProjectGroupsPage,
  ]
})
export class ProjectGroupsPageModule {}
