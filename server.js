var express = require('express');
var config = require('config');
var path = require('path');

var app = express();

// Middleware
app.use(express.static(path.join(__dirname, 'app/www')));

app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(require('express-session')({
  secret: config.get('session.secret'),
  // cookie: { maxAge: 60000 },
  resave: false,
  saveUninitialized: false
}));

app.disable('x-powered-by')

app.use(function (req, res, next) {
  res.locals.req = req;
  next();
});

// Routes
var mainRouter = require('./routes/main');

app.use(mainRouter);

app.get('*', function(req, res) {
    res.redirect('/');
});

const port = config.get('console.port');
app.listen(port, function() {});


if(process.argv[2] == 'dev') {
  var gulp = require('gulp');
  require('./gulpfile');
  if (gulp.tasks.styles) gulp.start('styles');
  if (gulp.tasks.scripts) gulp.start('scripts');
  if (gulp.tasks.watch) gulp.start('watch');
}
