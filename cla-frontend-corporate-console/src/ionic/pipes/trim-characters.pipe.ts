import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'trimCharacters'
})
export class TrimCharactersPipe implements PipeTransform {
    transform(v: string, shortenBy: number): string {
        if (v !== undefined && v !== null) {
            return v ? (v.length > shortenBy ? v.slice(0, shortenBy).concat('...') : v) : '';
        }
        return '';
    }
}
