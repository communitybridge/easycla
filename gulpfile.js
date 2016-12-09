var gulp = require('gulp');
var minifyCSS = require('gulp-minify-css');
var uglify = require('gulp-uglify');
var concat = require('gulp-concat');

gulp.task('styles', function(){
  return gulp.src('public/assets/src/css/*.css')
  .pipe(concat('all.min.css'))
  .pipe(minifyCSS({keepSpecialComments: 1}))
  .pipe(gulp.dest('public/assets/dist/'));
});

gulp.task('scripts', function(){
  return gulp.src('public/assets/src/js/*.js')
  .pipe(concat('all.min.js'))
  .pipe(uglify())
  .pipe(gulp.dest('public/assets/dist/'));
});

gulp.watch('watch', function(){
  gulp.watch('public/assets/src/css/*.js', ['scripts']);
});
