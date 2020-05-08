import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
    name: 'localTimeZone'
})
export class LocalTimeZonePipe implements PipeTransform {
    transform(utc: string): string {
        if (utc) {
            return new Date(utc).toLocaleString();
        }
        return '-';
    }
}
