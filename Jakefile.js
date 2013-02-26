var fs = require('fs'),
    eyes = require('eyes'),
    xml2js = require('xml2js'),
    moment = require('moment'),
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

desc('blog post on home page');
task('index-entry', function (params) {
  console.log(index_entry('my awesome title', '12 Feb 2031'));
});

desc('blog post');
task('post', function (params) {
  console.log(post('my awesome title', '12 Feb 2031', 'blog post!'));
});

desc('Generate all wp blog posts.');
task('wordpress', function (params) {
  var parser = new xml2js.Parser();
  jake.rmRf(__dirname + '/2007');
  jake.rmRf(__dirname + '/2008');
  jake.rmRf(__dirname + '/2009');
  jake.rmRf(__dirname + '/2010');
  jake.rmRf(__dirname + '/2011');
  jake.rmRf(__dirname + '/2012');
  jake.rmRf(__dirname + '/2013');

  parser.on('end', function(result) {
    //eyes.inspect(result.rss.channel[0].item);

    var to_inspect = _.filter(result.rss.channel[0].item, function(arr) {
      //console.log(arr['wp:status']);
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
      console.log(entry.title[0]);
      var content = entry["content:encoded"][0];
      var date = entry.pubDate[0];
      var asDate = moment(date, 'ddd, DD MMM YYYY HH:mm:ss Z');
      var dateString = asDate.utc().format('MMMM D, YYYY');
      console.log(dateString);
      content = content.replace(/(\r\n|\n|\r)/gm, "<br/>");
      var post_name = entry["wp:post_name"][0];

      var year = asDate.utc().year();
      var month = asDate.utc().format('MM');
      var day = asDate.utc().format('DD');

      var blog = post(entry.title[0], dateString, content);

      var loc = year + '/' + month + '/' + day + '/' + post_name;
      jake.mkdirP(loc);

      var x = i++;
      
      index_entries += index_entry(entry.title[0], dateString, loc + '/index.html');
      fs.appendFile(__dirname + '/' + loc + '/' + 'index.html', blog, function (err) {
        if (err) throw err;
        console.log('wrote blog '+ post_name);
      });
    });
    jake.rmRf(__dirname + '/index.html');
    fs.appendFile(__dirname + '/index.html', index(index_entries), function (err) {
      if (err) throw err;
      console.log('wrote index');
    });
  });

  parser.on('error', function(result) {
    console.log(result);
  });

  fs.readFile(__dirname + '/blog.xml', function(err, data) {
    parser.parseString(data);
  });
});
