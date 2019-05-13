import {Component, ViewChild} from "@angular/core";
import {
  NavController,
  NavParams,
  IonicPage, Nav, Events, AlertController
} from "ionic-angular";
import {ClaService} from "../../../services/cla.service";
import {Restricted} from "../../../decorators/restricted";

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
  selectedTemplate = {};
  templateValues = {};
  step = 'selection';

  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public claService: ClaService,
  ) {
    this.sfdcProjectId = navParams.get("projectId");
    this.projectTemplateId = navParams.get("projectTemplateId");
    this.getDefaults();
  }

  getDefaults() {
    this.getTemplates();
  }

  getTemplates () {
    this.claService.getTemplates().subscribe(templates => {
      this.templates = templates;
    });
  }

  ngOnInit() {
  }

  selectTemplate (template) {
    this.selectedTemplate = template;
    this.step = 'values';
  }

  reviewSelectedTemplate () {

  }

  goToFirstStep () {
    this.step = 'selection';
  }

  backToProject() {
    this.navCtrl.pop();
  }
}
