-Guice's toProvider Sometimes Considered Harmful
July 21, 2008
Perfection is what I wanted. Writing code examples for [a book](https://www.apress.com/gb/book/9781590599976) is harder than you think; Not only does the code need to be right, it also needs to be short and lightweight. Don't involve API's you don't need, don't create types if you can do without them, that sort of thing.

So I was nearing the handoff deadline and I was going over my code examples. I corrected some and tried to shorten others. For example, in the Chapter 6 example I saw the following code, which I used to abstract session handling from my view logic*:

```
public class UserTokenProvider implements Provider<UserToken> {
  @Inject private HttpSession session;
  public UserToken get() {
    synchronized (session) {
      return (UserToken) session.getAttribute(UserToken.KEY);
    }
  }
}
```

```
public class WebModule extends AbstractModule {
  protected void configure() {
    bind(UserToken.class).toProvider(UserTokenProvider.class);
  }
}
```

Surely I wasn't going to list two classes just for a simple provider, so I quickly decided to get rid of the extra class as follows:

```
public class WebModule extends AbstractModule {
  protected void configure() {
    bind(UserToken.class).toProvider(new Provider<UserToken>() {
      @Inject private HttpSession session;
      public UserToken get() {
        synchronized (session) {
          return (UserToken)session.getAttribute(UserToken.KEY);
        }
      }
    }); // no scope!
  }
}
```

Now, months later it struck me. This is wrong! (and arguably even less readable) In the shortened version I am now giving Guice a provider *instance* instead of a provider *type*. Using that `toProvider` overload is comparable to using Guice's `toInstance`: it will just reuse that instance for all requests to that `Key`, **disregarding** **all scopes**. And when I say all scopes, I mean all scopes. **It also ignores the default "no scope"**. Using the default scope, Guice will create an instance each time that `Key` gets requested. However, if you bypass all scopes with `toInstance` or the `toProvider` instance overload, Guice will simply reuse your instance.

In my example I was depending on the fact that the `HttpSession` injected in the provider was always going to be the right one. In the original example it worked as expected. I gave Guice a provider type, and used the default "no scope" so that a new instance would get created for each incoming HTTP request. Guice would inject the right `HttpSession` instance depending on the request.

Because the second example uses `toProvider` with an instance instead of a type, it behaves radically different. Ignoring "no scope", Guice will just inject the first `HttpSession` instance it finds and from then on it will leave that provider instance alone (a scope widening injection, if you will). Any subsequent requests, from possibly different sessions, will reuse that instance, leading to session corruption and a significant security risk.

Needless to say, I have fixed my code example and published the update to Apress.

Save yourself from such an embarrassment: **Remember that toInstance and the toProvider instance overload ignore all scopes, including "no scope". Avoid using them: use asEagerSingleton to load instances eagerly.** Only use these short cuts to fit code on a slide, or in the case of `toInstance`, to bind constants with better type safety.

\* I am not using Guice's session scope because the `UserToken` is used as an authentication token in the session. Using scopes, Guice might end up creating a security token for you when it gets requested, which is obviously not what you'd want.
