-Silverlight: an offer you can't refuse
May 4, 2007
If you haven't heard of it yet, [Silverlight](https://www.microsoft.com/silverlight/) (formerly WPF/E) is the RIA platform Microsoft announced recently. But at the [Mix 07](http://www.visitmix.com/) conference this week they surprised many: Silverlight 1.1 will ship with a stripped-down [CLR](http://en.wikipedia.org/wiki/Common_Language_Runtime). At the same time, they also announced the Dynamic Language Runtime (DLR), a framework for implementing dynamic languages on top of the CLR. In short, all this enables you to target the browser in your favorite .NET targeted language: C#, VB, Python, Ruby and even an ECMAScript 3 implementation  that runs on the CLR (not to be confused with JScript.NET, which has different goals). Anyway, 1.0 beta and 1.1 alpha have shipped, so you can try it out yourself.

For people who never really liked Javascript, like me, this is great news. This means that you will no longer have to rely on Javascript to get some RIA/Ajax behaviour in your web app (built-in Javascript engine, or Adobe Flex and the like). Being able to use the language you know best in all tiers makes absolute sense.

Now if only Java had something like this (and God no, applets are NOT equivalent), because after all, Java is the language I know best, and damn right I would like to use it in the browser as well. An intriguing possibility is to use [IKVM](http://www.ikvm.net) for targeting the Silverlight CLR. That would still be a little bit involved, but it could work.

Sun should wake up and smell the Java. I want real Java in the browser. Not [GWT](http://code.google.com/webtoolkit/). Not Java Web Start. Applets done right. Friction-free deployment, and fast as lightning.

Thanks in advance.
