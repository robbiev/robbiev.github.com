Migrating URLs from WordPress to GitHub Pages
February 27, 2013
For years my blog has been hosted on [wordpress.com](http://garbagecollected.wordpress.com). As you have probably guessed from this blog entry's title I chose to host my blog on [GitHub Pages](http://pages.github.com) going forward.

I had the following goals for my blog migration:
- Generating WordPress-compatible URLs
- Optionally, having a migration path to a new URL naming scheme

The reason I'm even bothering to migrate old blog entries is that I want to be a good internet citizen. I always find it annoying when links break so I'll try to at least do my part. This means **migrating the content and generating URLs that are compatible** with WordPress. Luckily I was already using a custom domain name, making it all possible.

WordPress URLs look like this: `http://garbagecollected.org/2008/04/06/dollarmaps/`.

Targeting GitHub Pages you then have two options:
1. mkdir -p 2008/04/06; touch 2008/04/06/dollarmaps.html
2. mkdir -p 2008/04/06/dollarmaps; touch 2008/04/06/dollarmaps/index.html

The first option *appears* to work at first; GitHub [automatically finds](https://github.com/holman/feedback/issues/231) the '.html' file with the correct name and serves its contents. However, the show stopper for me is that this no longer works if you add a trailing slash, something which WordPress always seems to do. So [http://garbagecollected.org/2008/04/06/dollarmaps](http://garbagecollected.org/2008/04/06/dollarmaps) works, but [http://garbagecollected.org/2008/04/06/dollarmaps/](http://garbagecollected.org/2008/04/06/dollarmaps/) doesn't. Not acceptable. This means that the only other option is to **create directories with an index.html**.

Now what about changing the URL naming scheme? Well, it turns out this is not really possible using GitHub pages. You would normally issue an HTTP redirect, but that's not supported. So then you can either generate [placeholder pages](http://stackoverflow.com/questions/10178304/github-jekyll-old-pages-redirection-best-approach) that redirect using JavaScript or the &lt;meta&gt; tag, or create a [custom 404 page](http://stackoverflow.com/questions/10178304/github-jekyll-old-pages-redirection-best-approach) that does the same. These are all more of a hack; you'll need to start [telling Google](http://support.google.com/webmasters/bin/answer.py?hl=en&answer=139394) that there is really only one page and things like that. For now I decided that new blog entries will follow the WordPress naming scheme until I find an easier hosting solution than GitHub's Pages.
