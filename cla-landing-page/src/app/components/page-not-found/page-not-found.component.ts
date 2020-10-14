import { Component, OnInit } from '@angular/core';
import { AppSettings } from 'src/app/config/app-settings';
import { REDIRECT_AUTH_ROUTE } from 'src/app/config/auth-utils';
import { EnvConfig } from 'src/app/config/cla-env-utils';
import { AuthService } from 'src/app/core/services/auth.service';
import { StorageService } from 'src/app/core/services/storage.service';

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
  styleUrls: ['./page-not-found.component.scss']
})
export class PageNotFoundComponent implements OnInit {
  message: string;
  constructor(
    private authService: AuthService,
    private storageService: StorageService
  ) { }

  ngOnInit(): void {
    if (this.authService.loggedIn) {
      const type = JSON.parse(this.storageService.getItem('type'));
      const redirectConsole = (type === 'Projects') ? AppSettings.PROJECT_CONSOLE_LINK : AppSettings.CORPORATE_CONSOLE_LINK;
      window.open(EnvConfig.default[redirectConsole] + REDIRECT_AUTH_ROUTE, '_self');
    } else {
      this.message = 'The page you are looking for was not found.';
    }
  }
}
