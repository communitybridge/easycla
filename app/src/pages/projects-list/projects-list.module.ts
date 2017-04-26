import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectsListPage } from './projects-list';

@NgModule({
  declarations: [
    ProjectsListPage
  ],
  imports: [
    IonicPageModule.forChild(ProjectsListPage)
  ],
  entryComponents: [
    ProjectsListPage
  ]
})
export class ProjectsListPageModule {}
