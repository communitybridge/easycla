import { NgModule, ErrorHandler } from '@angular/core';
import { IonicApp, IonicModule, IonicErrorHandler } from 'ionic-angular';
import { MyApp } from './app.component';
import { ProjectsListPage } from '../pages/projects-list/projects-list';
import { AddProjectPage } from '../pages/add-project/add-project';
import { MemberPage } from '../pages/member/member';
import { ProjectPage } from '../pages/project/project';
import { ContactUpdate } from '../pages/contact-update/contact-update';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';
import { UploadButtonComponent } from '../components/upload-button/upload-button';

@NgModule({
  declarations: [
    MyApp,
    ProjectsListPage,
    AddProjectPage,
    MemberPage,
    ProjectPage,
    ContactUpdate,
    UploadButtonComponent,
  ],
  imports: [
    IonicModule.forRoot(MyApp)
  ],
  bootstrap: [IonicApp],
  entryComponents: [
    MyApp,
    ProjectsListPage,
    AddProjectPage,
    MemberPage,
    ProjectPage,
    ContactUpdate,
  ],
  providers: [
    StatusBar,
    SplashScreen,
    {provide: ErrorHandler, useClass: IonicErrorHandler}
  ]
})
export class AppModule {}
