var fs = require('fs'),
    xml2js = require('xml2js'),
    moment = require('moment'),
    marked = require('marked'),
    path = require('path'),
    $ = require('cheerio'),
    _ = require('underscore');

var baseLocation = '../';
var blogLocation = '/../../blog.xml';
var blogEntryLocation = baseLocation + 'blog_entries/';
var mdBlogEntryLocation = baseLocation + 'blog_entries_md/';
var indexLocation = baseLocation + 'index.html';
var dateFormat = 'MMMM D, YYYY';
var wpDateFormat = 'ddd, DD MMM YYYY HH:mm:ss Z';

var indexEntry = function (title, date, path) {
  var indexEntry = fs.readFileSync('index-entry-template.html').toString()
  var html = $.load(indexEntry);
  html('a').text(title);
  html('a').attr('href', path);
  html('.date').text(date);
  return html.html();
};

var index = function (entries) {
  var indexEntry = fs.readFileSync('index-template.html').toString()
  var html = $.load(indexEntry);
  html('.home').html(entries);
  return html.html();
};

var post = function (title, date, entry) {
  var indexEntry = fs.readFileSync('post-template.html').toString()
  var html = $.load(indexEntry);
  html('title').text(title);
  html('h1').text(title);
  html('.date').text(date);
  html('.entry').html(entry);
  return html.html();
};

desc('Generate all blog posts.');
task('default', function (params) {
  var indexEntries = [];

  var generateEntries = function(listing, contentProcessor, extension) {
    var entries = _.filter(listing, function(e) { 
      return fs.statSync(e).isFile();
    });
  
    var entriesByLine = _.map(entries, function(entry) { 
      var blogFile = fs.readFileSync(entry).toString();
      var blogAsArray = blogFile.split(/\n/);
      return { file: entry, splitFile: blogAsArray };
    });
  
    _.each(entriesByLine, function(e) {
      var entry = e.file;
      var blogAsArray = e.splitFile.slice(0); // clone
  
      // pop the first two lines: title, date
      var title = blogAsArray[0];
      var date  = blogAsArray[1];
      blogAsArray.splice(0, 2);
  
      // the rest is the blog content
      var blogContent = blogAsArray.join('\n');

      var processedBlogContent = contentProcessor(blogContent);
  
      // generate HTML for the title, date and content
      var content = post(title, date, processedBlogContent);
  
      // generate blog entry path
      var asDate = moment(date, dateFormat);
      var year = asDate.year();
      var month = asDate.format('MM');
      var day = asDate.format('DD');
      var loc = year + '/' + month + '/' + day + '/' + path.basename(entry, extension);
      console.log(loc);
  
      // write blog entry HTML to disk
      jake.rmRf(baseLocation+loc);
      jake.mkdirP(baseLocation+loc);
      var file = baseLocation + loc + '/index.html';
      fs.writeFileSync(file, content);
      console.log('wrote blog ' + title);
  
      // save this entry to be included on the home page
      indexEntries.push({ 
        html: indexEntry(title, date, loc + '/'), 
        timestamp: moment(e.splitFile[1], dateFormat).unix() 
      });
    });
  };

  // process HTML blog entries
  var listing = jake.readdirR(blogEntryLocation);
  generateEntries(listing, _.identity, '.html');

  // process Markdown blog entries
  var mdListing = jake.readdirR(mdBlogEntryLocation);
  generateEntries(mdListing, marked, '.md');

  // sort the index entries
  var sortedIndexEntries = _.sortBy(indexEntries, function(e) {
    return e.timestamp
  }).reverse();

  // generate index.html
  var indexHtml = index(_.pluck(sortedIndexEntries, 'html').join(''));

  // write index.html
  jake.rmRf(indexLocation);
  fs.writeFileSync(indexLocation, indexHtml);
});
