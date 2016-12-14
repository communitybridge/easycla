var gulp = require('gulp');
var minifyCSS = require('gulp-minify-css');
var uglify = require('gulp-uglify');
var concat = require('gulp-concat');

gulp.task('styles', function(){
  console.log('Concatenating and minifying CSS files from /public/assets/src/css to /public/assets/dist');
  return gulp.src('public/assets/src/css/*.css')
  .pipe(concat('all.min.css'))
  .pipe(minifyCSS({keepSpecialComments: 1}))
  .pipe(gulp.dest('public/assets/dist/'));
});

gulp.task('scripts', function(){
  console.log('Concatenating and minifying JS files from /public/assets/src/js to /public/assets/dist');
  return gulp.src('public/assets/src/js/*.js')
  .pipe(concat('all.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('public/assets/dist/'));
});

gulp.task('watch', function(){
  console.log('Starting Gulp Watch...');
  gulp.watch('public/assets/src/css/*.css', ['styles']);
  gulp.watch('public/assets/src/js/*.js', ['scripts']);
});

gulp.task('default', ['styles', 'scripts', 'watch']);
