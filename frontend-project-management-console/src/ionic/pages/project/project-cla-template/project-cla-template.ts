import {Component, ViewChild} from "@angular/core";
import {
  NavController,
  ModalController,
  NavParams,
  IonicPage, Nav, Events, AlertController
} from "ionic-angular";
import { CincoService } from "../../../services/cinco.service";
import { KeycloakService } from "../../../services/keycloak/keycloak.service";
import { SortService } from "../../../services/sort.service";
import { ClaService } from "../../../services/cla.service";
import { RolesService } from "../../../services/roles.service";
import { Restricted } from "../../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated", "isPmcUser"]
})
@IonicPage({
  segment: "project/:projectId/cla/template/:projectTemplateId"
})
@Component({
  selector: "project-cla-template",
  templateUrl: "project-cla-template.html"
})
export class ProjectClaTemplatePage {
  loading: any;
  sfdcProjectId: string;
  projectTemplateId: string;
  templates: any[] = [];

  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
    public alertCtrl: AlertController,
    public claService: ClaService,
    public rolesService: RolesService,
    public events: Events
  ) {
    this.sfdcProjectId = navParams.get("projectId");
    this.projectTemplateId = navParams.get("projectTemplateId");
    this.getDefaults();
  }

  getDefaults() {
    this.getTemplates();
  }

  getTemplates () {
    this.claService
      .getTemplates()
      .subscribe(templates => {
        this.templates = templates
      })
  }

  ngOnInit() {
  }

  selectTemplate (template) {

    this.backToProject();
  }

  backToProject() {
    this.navCtrl.pop();
  }
}
