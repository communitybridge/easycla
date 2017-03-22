if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var config = require('config');
var CasStrategy = require('passport-cas').Strategy;
var path = require('path');
var flash = require('connect-flash');
var url = require('url');

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

app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(require('express-session')({
  secret: config.get('session.secret'),
  // cookie: { maxAge: 60000 },
  resave: false,
  saveUninitialized: false
}));
app.use(flash());

app.use(passport.initialize());
app.use(passport.session());

app.disable('x-powered-by')

app.use(function (req, res, next) {
  res.locals.req = req;
  next();
});

// Routes
var mainRouter = require('./routes/main');
var adminRouter = require('./routes/admin');
var projectsRouter = require('./routes/projects');
var membersRouter = require('./routes/members');
var mailingRouter = require('./routes/mailing');
var aliasesRouter = require('./routes/aliases');

app.use(mainRouter);
app.use(adminRouter);
app.use(projectsRouter);
app.use(membersRouter);
app.use(mailingRouter);
app.use(aliasesRouter);

app.get('*', function(req, res) {
    res.redirect('/');
});

// AWS  nginx proxy server uses 8081 by default
const appPort = 8081;
app.listen(appPort, function() {});
console.log("Node App listening port: ", appPort);

// Docker resolves the port mapping to "console.endpoint.port:8081"
const pmcURL = config.get('console.endpoint');
console.log('Docker project-management-console-url: ' + pmcURL);

if(process.argv[2] == 'dev') {
  var gulp = require('gulp');
  require('./gulpfile');
  if (gulp.tasks.styles) gulp.start('styles');
  if (gulp.tasks.scripts) gulp.start('scripts');
  if (gulp.tasks.watch) gulp.start('watch');
}

passport.use(new CasStrategy({
  version: 'CAS3.0',
  validateURL: '/serviceValidate',
  ssoBaseURL: 'https://identity.linuxfoundation.org/cas',
  serverBaseURL: pmcURL
}, function(login, done) {
  return done(null, login);
}));

passport.serializeUser(function(user, callback) {
  callback(null, user);
});

passport.deserializeUser(function(obj, callback) {
  callback(null, obj);
});
