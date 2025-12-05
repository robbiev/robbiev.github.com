-Press F3 to continue: initial thoughts on the JavaFX RIA platform
May 8, 2007
It all started with [this](http://www.infoq.com/news/2007/05/javafx-script) InfoQ article. Apparently, Sun will announce an early version (alpha?) of <a href="http://www.sun.com/software/javafx/">**JavaFX**</a> at [JavaOne](http://sun.com/javaone), an open source RIA platform that uses the [F3](http://blogs.sun.com/chrisoliver/resource/f3.html) scripting language and has a set of UI libraries and a programming model that would rival those of [Silverlight](https://www.microsoft.com/silverlight/) or [Flex](https://flex.apache.org/). The [F3 programming language](http://blogs.sun.com/chrisoliver/resource/f3.html) itself feels a bit like a mix of [Groovy](http://www.groovy-lang.org/) and Javascript 2 to me, but is looks like a powerful solution.

Of course it remains to be seen if it was really needed to have yet another programming language. Scrolling to the F3 spec, I can see one thing that they didn't get exactly right: integrated query syntax. I've seen some interviews and presentations on Microsoft's LINQ project recently, including the one Anders Hejlsberg gave at Mix 07. Someone in the audience asked "why do you guys put the from clause before the select clause?". The answer was simple: how would the IDE know how to help you if you start with the select clause? Select it from what? Subtle, but important. But it looks like a nice language, and the type inference thing is what I like to call "dynamic typing done right". And it has excellent Java interoperability, so there's a whole universe of libraries available to you. Check out [Chris Oliver's blog](http://blogs.sun.com/chrisoliver/category/F3) for more info on F3.

Although JavaFX looks like a worthy addition to the Java platform, they'll have to get some important things right:

- Tooling support: GUI builder, IDE plugins for F3, ...
- Excellent browser support: easy deployment and fast (Silverlight and Flex got this one right)

I fear that Sun, as the "spec company", will only implement a part of that story, and will count on the community to complete the picture. Well, we all know that some things just need some company throwing some money at it to get the whole story. Including timing...

I guess I'm exited. Let's wait for the JavaOne announcement before we draw any conclusions (based on some rumors). But draw 'em we will ;-)

**Update:** [this blog](https://edtechdev.wordpress.com/2007/05/08/f3-is-now-javafx/) links to some interesting stuff, including some hands-on exercises. Looks like they're using Java Web Start for deployment. Where's the browser integration?

**Update:** added link to official JavaFX page, lots of info there. The IDE support is there, but they're using Java's mechanism for deployment (standalone, Web Start, Applet).
