import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app.module';
import { enableProdMode } from '@angular/core';

import { KeycloakService } from '../services/keycloak/keycloak.service';

enableProdMode();

//KeycloakService.init({ onLoad: 'check-sso', checkLoginIframeInterval: 1 })
KeycloakService.init({ onLoad: 'login-required' })
    .then(() => {
        platformBrowserDynamic().bootstrapModule(AppModule);
    })
    .catch((e: string) => {
        console.log('Error in ng4 bootstrap: ' + e);
    });
