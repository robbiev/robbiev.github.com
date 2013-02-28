var fs = require('fs'),
    xml2js = require('xml2js'),
    moment = require('moment'),
    path = require('path'),
    $ = require('cheerio'),
    _ = require('underscore');

var baseLocation = '../';
var blogLocation = '/../../blog.xml';
var blogEntryLocation = baseLocation + 'blog_entries/';
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
  var indexEntries = '';
  var entries = jake.readdirR(blogEntryLocation);

  entries = _.filter(entries, function(e) { 
    return fs.statSync(e).isFile();
  });

  entries = _.map(entries, function(entry) { 
    var blogFile = fs.readFileSync(entry).toString();
    var blogAsArray = blogFile.split(/\n/);
    return { file: entry, splitFile: blogAsArray };
  });

  entries = _.sortBy(entries, function(e) {
    return moment(e.splitFile[1], dateFormat).unix();
  }).reverse();

  _.each(entries, function(e) {
    var entry = e.file;
    var blogAsArray = e.splitFile;

    var title = blogAsArray[0];
    var date  = blogAsArray[1];
    blogAsArray.splice(0, 2);

    var blogContent = blogAsArray.join('\n');
    var content = post(title, date, blogContent);

    var asDate = moment(date, dateFormat);
    var year = asDate.year();
    var month = asDate.format('MM');
    var day = asDate.format('DD');
    var loc = year + '/' + month + '/' + day + '/' + path.basename(entry, '.html');
    console.log(loc);

    jake.rmRf(baseLocation+loc);
    jake.mkdirP(baseLocation+loc);
    var file = baseLocation + loc + '/index.html';
    fs.writeFileSync(file, content);
    console.log('wrote blog ' + title);

    indexEntries += indexEntry(title, date, loc + '/');
  });
  var indexHtml = index(indexEntries);
  jake.rmRf(indexLocation);
  fs.writeFileSync(indexLocation, indexHtml);
});

desc('Generate blog entries from WordPress XML.');
task('wp', function (params) {
  var parser = new xml2js.Parser();

  parser.on('end', function(result) {
    var to_inspect = _.filter(result.rss.channel[0].item, function(arr) {
      return arr['wp:status'][0] === 'publish';
    });

    var i = 0;
    var index_entries = "";
    _.each(to_inspect, function(entry) {
      var date = entry.pubDate[0];
      var asDate = moment(date, wpDateFormat);
      var dateString = asDate.utc().format(dateFormat);

      var content = entry["content:encoded"][0];
      content = content.replace(/(\r\n|\n|\r)/gm, "<br/>");
      var to_write = entry.title[0] + '\n' + dateString + '\n' + content;

      var post_name = entry["wp:post_name"][0];
      var fileToWrite = blogEntryLocation + post_name + '.html';

      jake.rmRf(fileToWrite);
      fs.writeFileSync(fileToWrite, to_write);
    });
  });

  parser.on('error', function(result) {
    console.log(result);
  });

  fs.readFile(__dirname + blogLocation, function(err, data) {
    parser.parseString(data);
  });
});
