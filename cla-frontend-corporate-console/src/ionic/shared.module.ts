import { NgModule } from '@angular/core';
import { TrimCharactersPipe } from './pipes/trim-characters.pipe';
import { LocalTimeZonePipe } from './pipes/local-timezone.pipe';

@NgModule({
    declarations: [
        TrimCharactersPipe,
        LocalTimeZonePipe
    ],
    imports: [

    ],
    exports: [
        TrimCharactersPipe,
        LocalTimeZonePipe
    ],
    providers: []
})
export class SharedModule { }
