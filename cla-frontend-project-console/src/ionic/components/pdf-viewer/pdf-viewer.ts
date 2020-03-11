import { Component, Input } from '@angular/core';

@Component({
  selector: 'easy-pdf-viewer',
  templateUrl: 'pdf-viewer.html'
})
export class PdfViewerComponent {
  zoom: number = 0.7;
  @Input() src: string;

  constructor() { }

  //zoom in pdf
  increment() {
    if (this.zoom < 1) {
      this.zoom = this.zoom + 0.1;
    }
  }

  //zoom out pdf
  decrement() {
    if (this.zoom > 0) {
      this.zoom = this.zoom - 0.1;
    }
  }

  //download pdf
  download() {
    console.log(this.src);
    const downloadLink = document.createElement('a');
    downloadLink.style.display = 'none';
    document.body.appendChild(downloadLink);
    downloadLink.setAttribute('href', this.src);
    downloadLink.setAttribute('download', 'true');
    downloadLink.setAttribute('target', '_blank');
    downloadLink.click();
    document.body.removeChild(downloadLink);
  }
}
