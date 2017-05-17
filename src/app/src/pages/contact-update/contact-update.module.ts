import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ContactUpdate } from './contact-update';
import { ComponentsModule } from '../../components/components.modules';

@NgModule({
  declarations: [
    ContactUpdate,
  ],
  imports: [
    IonicPageModule.forChild(ContactUpdate),
    ComponentsModule
  ],
  entryComponents: [
    ContactUpdate,
  ]
})
export class ContactUpdateModule {}
