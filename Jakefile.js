var fs = require('fs'),
    eyes = require('eyes'),
    xml2js = require('xml2js'),
    moment = require('moment'),
    _ = require('underscore');

var parser = new xml2js.Parser();

desc('This is the default task.');
task('default', function (params) {

  jake.rmRf("out");
  jake.mkdirP("out");

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
    });

    var i = 0;
    _.each(sorted, function(entry) {
      console.log(entry.title[0]);
      var content = entry["content:encoded"][0];
      var date = entry.pubDate[0];
      var asDate = moment(date, 'ddd, DD MMM YYYY HH:mm:ss Z');
      console.log(asDate.utc().format('MMMM D, YYYY'));
      content = content.replace(/(\r\n|\n|\r)/gm, "<br/>");
      //console.log(content);
      var post_name = entry["wp:post_name"][0];

      var year = asDate.utc().year();
      var month = asDate.utc().format('MM');
      var day = asDate.utc().format('DD');

      var loc = 'out/'+ year + '/' + month + '/' + day + '/' + post_name;
      jake.mkdirP(loc);

      var x = i++;
      fs.appendFile(__dirname + '/' + loc + '/' + 'index.html', content, function (err) {
        if (err) throw err;
        console.log('wrote blog '+ post_name);
      });
    });
  });

  parser.on('error', function(result) {
    console.log(result);
  });

  fs.readFile(__dirname + '/blog.xml', function(err, data) {
    parser.parseString(data);
  });
});
