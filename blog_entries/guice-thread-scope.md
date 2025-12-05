-Guice Thread Scope
June 7, 2007
**UPDATE**: The reset thing turned out to be a bad idea, I created [simple Thread Scope](http://code.google.com/p/google-guice/issues/detail?id=114) instead. If you somehow seem to need a reset, consider [this](http://code.google.com/p/google-guice/wiki/Scopes).

Ever wanted to use [Guice](http://code.google.com/p/google-guice/)'s request scope outside of a web application? Try out [my thread scope implementation](http://code.google.com/p/google-guice/issues/attachment?aid=-4210064442097121541&name=GuiceThreadScope.txt) and let me know what you think. [Cast your vote](http://code.google.com/p/google-guice/issues/detail?id=100) if you want to see this scope added to Guice.

Here's an example. Create a `Module` that looks like this:

```
Injector i = Guice.createInjector(new Module() {
    public void configure(Binder binder) {
        binder.bindScope(ThreadScoped.class, CustomScopes.THREAD);
        binder.bind(ThreadCache.class).in(Scopes.SINGLETON);
        // add your custom classes
        binder.bind(SomeClass.class).in(CustomScopes.THREAD);
    }
});
```

Each thread is automatically initialized for the scope. To reset it, use the `ThreadCache`:

```
@Inject private ThreadCache threadCache
public void someMethod() {
    try {
        // do stuff in thread scope
    } finally {
        threadCache.reset();
    }
}
```

I'm thinking about loading this scope's initialization code lazily (so that you'll need to execute `threadCache.start()` or something like that). It currently doesn't matter for me and the memory overhead is low anyway, but what do you think?

PS: only use field injection for examples ;-)
