Lazy Loading For Mere Mortals
August 23, 2008
Easily the #1 mistake people make when using [Warp Persist](http://code.google.com/p/warp-persist) is they forget to start their `PersistenceService`, which is the API that abstracts the Hibernate `SessionFactory`. The plan for the next release is to start the `PersistenceService` lazily. Thus the first time a `SessionFactory` gets accessed, we will load it for you if you haven't called `PersistenceService.start()` yet. This is a simple change that will help you get your project going in no time.

The obvious choice for implementing this type of lazy loading as of Java 5 is [Double-Checked Locking](http://en.wikipedia.org/wiki/Double-checked_locking). But after coding it up for the Hibernate `SessionFactory`, I realized I would have to do the same for the JPA and DB4O support. Given the relative complexity of DCL, that kind of sucks.

But then I remembered something [Bob Lee](https://twitter.com/crazybob) said on [Twitter](https://www.twitter.com) the other day:

![Bob Lee on Twitter](/img/picture-1.png)

The man has a point.

I decided to roll our own utility class to lazily load object references. So I opened up the bible ([Effective Java, 2nd edition](http://java.sun.com/docs/books/effective/)) to look for the original pattern, and there it was. Page 283:

```
private volatile FieldType field;
FieldType getField() {
  FieldType result = field;
  if (result == null) {  // First check (no locking)
    synchronized(this) {
      result = field;
      if (result == null) // Second check (with locking)
        field = result = computeFieldValue();
    }
  }
  return result;
}
```

Scary. Now, there is no way we can reduce the size of that code much, but as it turns out we can make it a lot simpler to use. I called it [LazyReference](http://code.google.com/p/warp-persist/source/browse/trunk/warp-persist/src/com/wideplay/warp/util/LazyReference.java?r=128), here's an example usage:

```
private final LazyReference<SessionFactory> sessionFactory =
  LazyReference.of(new Provider<SessionFactory>() {
    public SessionFactory get() {
      // code to create SessionFactory
    }
});
```
To use this code, you just call `sessionFactory.get()`, and it will handle the lazy loading for you. Even if you are a Java master, there is no reason for you to repeat that complicated DCL code ever again. Enjoy!
