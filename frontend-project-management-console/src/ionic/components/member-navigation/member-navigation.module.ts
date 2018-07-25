import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { MemberNavigationComponent } from './member-navigation';

@NgModule({
  declarations: [
    MemberNavigationComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    MemberNavigationComponent,
  ]
})
export class MemberNavigationComponentModule {}
