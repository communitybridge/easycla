if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var bodyParser = require('body-parser');
var passport = require('passport');
var config = require('config');
var CasStrategy = require('passport-cas').Strategy;
var path = require('path');
var url = require('url');
var session = require('express-session');
var RedisStore = require('connect-redis')(session);
const UTIL = require('util');

var app = express();

// App config
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Middleware
app.use(express.static(path.join(__dirname, 'public')));

app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded
app.use(bodyParser.json()); // for parsing application/json
app.use(session({
  store: new RedisStore({
    host: config.get('console.redisHost'),
    port: config.get('console.redisPort')
  }),
  secret: config.get('session.secret'),
  resave: false,
  saveUninitialized: false
}));

app.use(passport.initialize());
app.use(passport.session());

app.disable('x-powered-by');

app.use(function (req, res, next) {
  res.locals.req = req;
  next();
});

var authMiddleware = function(req, res, next) {
  if (req.isAuthenticated()) return next();
  else return res.render('login');
}
app.use('/member-console', authMiddleware, express.static(path.join(__dirname, 'app/www')));

// Routes
var mainRouter = require('./routes/main');
var organizations = require('./routes/organizations');
var projectsRouter = require('./routes/projects');
var projectsMembersRouter = require('./routes/projects.members');
var projectsMembersContactsRouter = require('./routes/projects.members.contacts');
var organizationsContactsRouter = require('./routes/organizations.contacts');
var organizationsProjectsRouter = require('./routes/organizations.projects');
var mailingRouter = require('./routes/mailing');
var aliasesRouter = require('./routes/aliases');
var usersRouter = require('./routes/users');

app.use(mainRouter);
app.use(organizations);
app.use(projectsRouter);
app.use(projectsMembersRouter);
app.use(projectsMembersContactsRouter);
app.use(organizationsContactsRouter);
app.use(organizationsProjectsRouter);
app.use(mailingRouter);
app.use(aliasesRouter);
app.use(usersRouter);

app.get('*', function(req, res) {
    res.redirect('/');
});

// AWS  nginx proxy server uses 8081 by default
const APP_PORT = 8081;
app.listen(APP_PORT, function() {});
console.log("Node App listening port: ", APP_PORT);

// Docker resolves the port mapping to "console.endpoint.port:8081"
const MEMBER_CONSOLE_URL = config.get('console.endpoint');
console.log('Docker member-console-url: ' + MEMBER_CONSOLE_URL);

displayBanner();

passport.use(new CasStrategy({
  version: 'CAS3.0',
  validateURL: '/serviceValidate',
  ssoBaseURL: 'https://identity.linuxfoundation.org/cas',
  serverBaseURL: MEMBER_CONSOLE_URL
}, function(login, done) {
  return done(null, login);
}));

passport.serializeUser(function(user, callback) {
  // console.log(UTIL.inspect(user, false, null))
  callback(null, user.user);
});

passport.deserializeUser(function(obj, callback) {
  callback(null, obj);
});

function displayBanner() {
  let fs = require('fs')
  let max = 4;
  let min = 1;
  let version = Math.floor (Math.random() * (max - min + 1) ) + min;
  let banner = 'banners/' + version + '.txt';
  fs.readFile(banner, 'utf8', function(err, data) {
    if (err) throw err;
    console.log(data)
  });
};
