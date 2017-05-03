import { BrowserModule } from '@angular/platform-browser';
import { NgModule, ErrorHandler } from '@angular/core';
import { HttpModule } from '@angular/http';
import { IonicApp, IonicModule, IonicErrorHandler } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { MyApp } from './app.component';
import { ContactUpdate } from '../pages/contact-update/contact-update';
import { UploadButtonComponent } from '../components/upload-button/upload-button';
import { ActionPopover } from '../components/action-popover/action-popover';

@NgModule({
  declarations: [
    MyApp,
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
