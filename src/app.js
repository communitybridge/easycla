if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var config = require('config');
var CasStrategy = require('passport-cas').Strategy;
var path = require('path');
var flash = require('connect-flash');
var url = require('url');
var session = require('express-session');
var RedisStore = require('connect-redis')(session);
const util = require('util')

var app = express();

// App config
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Middleware
app.use(express.static(path.join(__dirname, 'public')));

app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(session({
  store: new RedisStore({
    host: config.get('console.redisHost'),
    port: config.get('console.redisPort')
  }),
  secret: config.get('session.secret'),
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

var authMiddleware = function(req, res, next) {
  if (req.isAuthenticated()) return next();
  else return res.render('login');
}
app.use('/pmc', authMiddleware, express.static(path.join(__dirname, 'app/www')));

// Routes
var mainRouter = require('./routes/main');
var adminRouter = require('./routes/admin');
var organizations = require('./routes/organizations');
var projectsRouter = require('./routes/projects');
var membersRouter = require('./routes/members');
var mailingRouter = require('./routes/mailing');
var aliasesRouter = require('./routes/aliases');
var usersRouter = require('./routes/users');

app.use(mainRouter);
app.use(adminRouter);
app.use(organizations);
app.use(projectsRouter);
app.use(membersRouter);
app.use(mailingRouter);
app.use(aliasesRouter);
app.use(usersRouter);

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

passport.use(new CasStrategy({
  version: 'CAS3.0',
  validateURL: '/serviceValidate',
  ssoBaseURL: 'https://identity.linuxfoundation.org/cas',
  serverBaseURL: pmcURL
}, function(login, done) {
  return done(null, login);
}));

passport.serializeUser(function(user, callback) {
  // console.log(util.inspect(user, false, null))
  callback(null, user.user);
});

passport.deserializeUser(function(obj, callback) {
  callback(null, obj);
});
