var fs = require('fs'),
    eyes = require('eyes'),
    xml2js = require('xml2js'),
    moment = require('moment'),
    path = require('path'),
    $ = require('cheerio'),
    _ = require('underscore');


var index_entry = function (title, date, path) {
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

desc('Generate all new blog posts.');
task('default', function (params) {
  var indexEntries = '';
  var entries = jake.readdirR('../blog_entries/');
  entries = _.filter(entries, function(e) { return fs.statSync(e).isFile(); });
  entries = _.map(entries, function(entry) { 
    var blogFile = fs.readFileSync(entry).toString();
    var blogAsArray = blogFile.split(/\n/);
    return { file: entry, splitFile: blogAsArray };
  });
  entries = _.sortBy(entries, function(e) { return moment(e.splitFile[1], 'MMM D, YYYY').unix(); }).reverse();
  _.each(entries, function(e) {
    var entry = e.file;
    var blogAsArray = e.splitFile;

    var title = blogAsArray[0];
    var date  = blogAsArray[1];
    blogAsArray.shift();
    blogAsArray.shift();
    var blogContent = blogAsArray.join('\n');
    var content = post(title, date, blogContent);
    var asDate = moment(date, 'MMM D, YYYY');
    var year = asDate.year();
    var month = asDate.format('MM');
    var day = asDate.format('DD');
    var loc = year + '/' + month + '/' + day + '/' + path.basename(entry, '.html');
    console.log(loc);
    jake.rmRf('../'+loc);
    jake.mkdirP('../'+loc);

    var file = '../'+loc + '/index.html';
    fs.appendFile(file, content, function (err) {
      if (err) throw err;
      console.log('wrote blog ' + title);
    });
    indexEntries += index_entry(title, date, loc);
  });
  var indexFile = fs.readFileSync('../index.html').toString();
  var html = $.load(indexFile);
  html('.home').prepend(indexEntries);
  fs.writeFileSync('../index.html', html.html());
});

desc('Generate all wordpress blog posts.');
task('wp', function (params) {
  var parser = new xml2js.Parser();
  jake.rmRf(__dirname + '/../2007');
  jake.rmRf(__dirname + '/../2008');
  jake.rmRf(__dirname + '/../2009');
  jake.rmRf(__dirname + '/../2010');
  jake.rmRf(__dirname + '/../2011');
  jake.rmRf(__dirname + '/../2012');
  jake.rmRf(__dirname + '/../2013');

  parser.on('end', function(result) {
    var to_inspect = _.filter(result.rss.channel[0].item, function(arr) {
      return arr['wp:status'][0] === 'publish';
    });

    var sorted = _.sortBy(to_inspect, function(entry) {
      var date = entry.pubDate[0];
      var asDate = moment(date, 'ddd, DD MMM YYYY HH:mm:ss Z');
      return asDate.utc().unix();
    }).reverse();

    var i = 0;
    var index_entries = "";
    _.each(sorted, function(entry) {
      var content = entry["content:encoded"][0];
      var date = entry.pubDate[0];
      var asDate = moment(date, 'ddd, DD MMM YYYY HH:mm:ss Z');
      var dateString = asDate.utc().format('MMMM D, YYYY');
      content = content.replace(/(\r\n|\n|\r)/gm, "<br/>");
      var post_name = entry["wp:post_name"][0];

      var year = asDate.utc().year();
      var month = asDate.utc().format('MM');
      var day = asDate.utc().format('DD');

      var blog = post(entry.title[0], dateString, content);

      var loc = year + '/' + month + '/' + day + '/' + post_name;
      jake.mkdirP('../'+loc);

      var x = i++;
      
      index_entries += index_entry(entry.title[0], dateString, loc + '/');
      fs.appendFile(__dirname + '/../' + loc + '/' + 'index.html', blog, function (err) {
        if (err) throw err;
        console.log('wrote blog '+ post_name);
      });
    });
    jake.rmRf(__dirname + '/../index.html');
    fs.appendFile(__dirname + '/../index.html', index(index_entries), function (err) {
      if (err) throw err;
      console.log('wrote index');
    });
  });

  parser.on('error', function(result) {
    console.log(result);
  });

  fs.readFile(__dirname + '/../../blog.xml', function(err, data) {
    parser.parseString(data);
  });
});
