var express = require('express');
var path = require('path');
var bodyParser = require('body-parser');

var app = express();

// App config
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Middleware
app.use(express.static(path.join(__dirname, 'public')));

// Modules required by Angular 2
app.use('/node_modules/zone.js/dist/', express.static(path.join(__dirname, 'node_modules/zone.js/dist/')));
app.use('/node_modules/reflect-metadata/', express.static(path.join(__dirname, 'node_modules/reflect-metadata/')));
app.use('/node_modules/systemjs/dist/', express.static(path.join(__dirname, 'node_modules/systemjs/dist/')));
app.use('/node_modules/core-js/client/', express.static(path.join(__dirname, 'node_modules/core-js/client/')));
app.use('/node_modules/@angular/', express.static(path.join(__dirname, 'node_modules/@angular/')));
app.use('/node_modules/angular2-in-memory-web-api/', express.static(path.join(__dirname, 'node_modules/angular2-in-memory-web-api/')));
app.use('/node_modules/rxjs/', express.static(path.join(__dirname, 'node_modules/rxjs/')));

app.use(require('morgan')('combined'));
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(require('express-session')({ secret: process.env['SESSION_SECRET'] != null ? process.env['SESSION_SECRET'] : 'lhb.sdu3erw lwfe rlfwe oThge3 823dwj34 @#kbdwe3 ghdklnj32lj l2303', resave: false, saveUninitialized: false }));

// Routes
var routes = require('./routes')
app.use(routes);

const port = process.env['UI_PORT'] != null ? process.env['UI_PORT'] : 8081
app.listen(port, function() {});
