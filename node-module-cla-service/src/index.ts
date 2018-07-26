import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClaService } from './cla.service';

export * from './cla.service';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [
  ],
  exports: [
  ]
})
export class ClaServiceModule {
  static forRoot(): ModuleWithProviders {
    return {
      ngModule: ClaServiceModule,
      providers: [ClaService]
    };
  }
}
