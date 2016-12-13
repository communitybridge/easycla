var newrelic = require('newrelic');
var express = require('express');
var passport = require('passport');
var CasStrategy = require('passport-cas').Strategy;
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

app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(require('express-session')({ secret: process.env['SESSION_SECRET'] != null ? process.env['SESSION_SECRET'] : 'lhb.sdu3erw lwfe rlfwe oThge3 825dwj34 @#kbdwe3 ghdklnj32lj l2303', resave: false, saveUninitialized: false }));

app.use(passport.initialize());
app.use(passport.session());

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

const port = process.env['UI_PORT'] != null ? process.env['UI_PORT'] : 8081
app.listen(port, function() {});

var serverBaseURL = process.env['CINCO_CONSOLE_URL'];
if(process.argv[2] == 'dev') {
  serverBaseURL = 'http://localhost:8081';
  var gulp = require('gulp');
  require('./gulpfile');
  if (gulp.tasks.styles) {
      console.log('Concatenating and minifying CSS files from /public/assets/src/css to /public/assets/dist');
      gulp.start('styles');
  }
  if (gulp.tasks.scripts) {
      console.log('Concatenating and minifying JS files from /public/assets/src/js to /public/assets/dist');
      gulp.start('scripts');
  }
}
if(!serverBaseURL.startsWith('http') ) serverBaseURL = 'http://' + serverBaseURL;
console.log('serverBaseURL: ' + serverBaseURL);

passport.use(new CasStrategy({
  version: 'CAS3.0',
  validateURL: '/serviceValidate',
  ssoBaseURL: 'https://identity.linuxfoundation.org/cas',
  serverBaseURL: serverBaseURL
}, function(login, done) {
  return done(null, login);
}));

passport.serializeUser(function(user, callback) {
  callback(null, user);
});

passport.deserializeUser(function(obj, callback) {
  callback(null, obj);
});
