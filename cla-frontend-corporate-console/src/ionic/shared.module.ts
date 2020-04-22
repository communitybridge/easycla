import { NgModule } from '@angular/core';
import { TrimCharactersPipe } from './pipes/trim-characters.pipe';

@NgModule({
    declarations: [
        TrimCharactersPipe
    ],
    imports:[

    ],
    exports: [
        TrimCharactersPipe
    ],
    providers: []
})
export class SharedModule { }
