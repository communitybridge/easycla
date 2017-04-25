import { BrowserModule } from '@angular/platform-browser';
import { NgModule, ErrorHandler } from '@angular/core';
import { HttpModule } from '@angular/http';
import { IonicApp, IonicModule, IonicErrorHandler } from 'ionic-angular';
import { MyApp } from './app.component';
// import { ProjectsListPage } from '../pages/projects-list/projects-list';
// import { MemberPage } from '../pages/member/member';
// import { ProjectPage } from '../pages/project/project';
import { ContactUpdate } from '../pages/contact-update/contact-update';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';
import { UploadButtonComponent } from '../components/upload-button/upload-button';
import { ActionPopover } from '../components/action-popover/action-popover';

@NgModule({
  declarations: [
    MyApp,
    // ProjectsListPage,
    // MemberPage,
    // ProjectPage,
    ContactUpdate,
    UploadButtonComponent,
    ActionPopover,
  ],
  imports: [
    BrowserModule,
    HttpModule,
    IonicModule.forRoot(MyApp),
  ],
  bootstrap: [IonicApp],
  entryComponents: [
    MyApp,
    // ProjectsListPage,
    // MemberPage,
    // ProjectPage,
    ContactUpdate,
    ActionPopover,
  ],
  providers: [
    StatusBar,
    SplashScreen,
    {provide: ErrorHandler, useClass: IonicErrorHandler}
  ]
})
export class AppModule {}
