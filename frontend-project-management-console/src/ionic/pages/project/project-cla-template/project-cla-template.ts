import {Component, ViewChild} from "@angular/core";
import {
  NavController,
  NavParams,
  IonicPage, Nav, Events, AlertController
} from "ionic-angular";
import {ClaService} from "../../../services/cla.service";
import {Restricted} from "../../../decorators/restricted";
import { DomSanitizer } from '@angular/platform-browser';

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
  projectId: string;
  templates: any[] = [];
  selectedTemplate: any;
  templateValues = {};
  pdfPath = {
    corporatePDFURL: '',
    individualPDFURL: ''
  };
  currentPDF = 'corporatePDFURL';
  step = 'selection';

  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public claService: ClaService,
    public sanitizer: DomSanitizer
  ) {
    this.sfdcProjectId = navParams.get("sfdcProjectId");
    this.projectId = navParams.get("projectId");
    this.getDefaults();
  }

  getDefaults() {
    this.getTemplates();
  }

  getTemplates() {
    this.claService.getTemplates().subscribe(templates => this.templates = templates);
  }

  ngOnInit() {
  }

  getPdfPath() {
    return this.sanitizer.bypassSecurityTrustResourceUrl(this.pdfPath[this.currentPDF]);
  }

  showPDF (type) {
    this.currentPDF = type;
  }

  selectTemplate(template) {
    this.selectedTemplate = template;
    this.step = 'values';
  }

  reviewSelectedTemplate() {

    var metaFields = this.selectedTemplate.metaFields; 
    metaFields.forEach(metaField => {
      if ( this.templateValues.hasOwnProperty(metaField.TemplateVariable)) {
        metaField.Value = this.templateValues[metaField.TemplateVariable]
      }

    });
    let data = {
      templateID: this.selectedTemplate.ID,
      metaFields: metaFields
    }

    this.claService.postClaGroupTemplate(this.projectId,  data)
      .subscribe(response => {
        this.pdfPath = response;
        console.log(this.pdfPath)
        this.goToStep('review');
      })
  }

  goToStep(step) {
    this.step = step;
  }

  backToProject() {
    this.navCtrl.pop();
  }

  viewDocusign() {}

  submit() {}
}
