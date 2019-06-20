// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Directive, ElementRef, Renderer2, Input, OnChanges, SimpleChange } from '@angular/core';

/**
 * Generated class for the LoadingDisplayDirective directive.
 *
 * See https://angular.io/docs/ts/latest/api/core/index/DirectiveMetadata-class.html
 * for more info on Angular Directives.
 */
@Directive({
  selector: '[loading-display]' // Attribute selector
})
export class LoadingDisplayDirective implements OnChanges {

  @Input('loading-display') loadingDisplay: any;

  constructor(
    public element: ElementRef,
    public renderer: Renderer2,
  ) {
  }
  ngOnInit() {
    this.renderer.addClass(this.element.nativeElement, 'loading-display-initial');
  }

  ngOnChanges(changes: {[propertyName: string]: SimpleChange}) {
    if (changes['loadingDisplay'] && !this.loadingDisplay) {
      this.renderer.addClass(this.element.nativeElement, 'loading-display-loaded');
    }
  }

}
