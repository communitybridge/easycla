import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AuthPage } from './auth';

@NgModule({
  declarations: [
    AuthPage,
  ],
  imports: [
    IonicPageModule.forChild(AuthPage),
  ],
})
export class AuthPageModule {}
