Non-functional Software
May 18, 2016

When we plan to build a software system we often discuss two sets of requirements:

1. **[Functional](https://en.wikipedia.org/wiki/Functional_requirement)**: *"When I do X, the software does Y."*.
1. **[Non-functional](https://en.wikipedia.org/wiki/Non-functional_requirement)**: privacy, performance, security, compatibility, ...

Recently I found [another definition for non-functional requirements](http://www.slideshare.net/littleidea/architecture-what-does-it-even-mean):

> requirements which, if not met, will make a system non-functional

Now, lately I've been thinking about this definition and how it relates to software dependencies. **Unless you're writing an operating system, your software always depends on other software**. At the very least you'll be making some [system calls](https://en.wikipedia.org/wiki/System_call) into the operating system. Using threads, sockets, files? These **operating system dependencies are a part of your non-functional requirements**.

What I'm getting at is that even if your library does not depend on other libraries (in the [NPM](https://www.npmjs.com/) or [Maven](https://maven.apache.org/) sense), *you still have dependencies* that are often implicit to the environment and you should consider documenting them. If your software's documentation says that *"When I do X, it does Y"* but it doesn't say anything about creating 100,000 threads to do Y, then you may want to consider documenting that. **Your non-functional requirement (or lack thereof) is my non-functional behaviour**.

Usually, however, the problem is much more subtle. A pet peeve of mine is that libraries often don't document **what thread a piece of code will run on**. What's wrong with you people? If you write a framework or your library accepts a callback function, in other words you control the entry point, then you absolutely need to document this.

Similarly, when is something blocking and when is it non-blocking? Is it okay to be blocked for long? If not, what do I do instead? If it accesses a networked service, does it automatically reconnect? With back-off? These are all things I need to know to ensure luck is not a factor in successfully running your software.

We should have a name for this kind of software - **non-functional software**:
> software that, by not documenting its non-functional behaviour, is non-functional
