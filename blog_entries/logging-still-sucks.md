Logging Still Sucks
March 7, 2013
It's 2013, and for some reason logging is still a pain in the ass when you use Java ([SLF4J](http://www.slf4j.org)) or .NET ([log4net](http://logging.apache.org/log4net/)).

SLF4J is indeed a better version of
[JCL](http://commons.apache.org/proper/commons-logging/). Using an
abstraction library is indeed better than just choosing a logging framework
outright. I can agree with that.  But in my opinion there is still something
wrong with the programming model:

- **Your `DEBUG` is not my `DEBUG`:** choosing a log level is highly subjective. Most
  projects don't have a standard practice for this, so everything just becomes a
  dumping ground. All you can say is that `DEBUG` will probably be more verbose
  than `INFO`. As a wise man once said: enabling log levels is like a [box of
    chocolates](http://www.imdb.com/title/tt0109830/quotes?qt=qt0373657).
- **Designed to be an afterthought:** it's too easy to ignore SLF4J and then start
  using it when you're already in production. Logging is actually a feature of the
  library you are using, but nobody will even realise it's there until they have
  problem. Problem in production?  Let's turn up logging! However, this can have
  side effects in terms of triggering suble bugs and degrading performance. Now
  you have two problems in production.
- **Simple case is no longer simple:** FFS, just let me dump this thing to sysout so I
  can move on with my life.

See the trend here? Only the API consumer knows what your `DEBUG` means for their
`DEBUG`. Also they need to know logging could be happening in the first place.
They need to take the time and map your logging to their logging (if they want
to) and make sure that everything still works as expected with logging enabled.
So as a API author, what do you do? You write an API.

```
public interface FooLogger {
    void FYI(Event event);
    void WTF(Event event);
    void OMG(Event event);
  }
```

So then as an API consumer, I implement FooLogger and give it to the API:

```
class MyFooLogger implements FooLogger {
    ...

    @Override
    public void FYI(Event event) {
      logger.INFO(event);
    }
    @Override
    public void WTF(Event event) {
      logger.WARN(event);
    }
    @Override
    public void OMG(Event event) {
      logger.ERROR(event);
    }
  }

  Library lib = Library.Create(new MyFooLogger());
  lib.doStuff();
```

An example of this style can actually be found in the [Apache Velocity](http://velocity.apache.org/engine/devel/developer-guide.html#Configuring_Logging) project.
For some reason [the
  SLF4J guys don't like this](http://www.slf4j.org/faq.html#optional_dependency):

>
  It is reasonable to assume that in most projects Wombat will be one dependency
  among many. If each library had its own logging wrapper, then each wrapper would
  presumably need to be configured separately. Thus, instead of having to deal
  with one logging framework, namely SLF4J, the user of Wombat would have to
  detail with Wombat's logging wrapper as well. The problem will be compounded by
  each framework that comes up with its own wrapper in order to make SLF4J
  optional. (Configuring or dealing with the intricacies of five different logging
  wrappers is not exactly exciting nor endearing.)

I don't see the problem. You write the connector code once for the lifetime of a
project. This code could last for years. The code will be simple as it will be
specific to your use case (e.g. log to an existing NoSQL system). You choose
when you do it. And seriously, if the framework makes this difficult then you
are essentially admitting their API sucks. Nothing more.
