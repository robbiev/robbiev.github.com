-Builder Pattern Deluxe
July 12, 2007
**Update:** code available at [http://code.google.com/p/garbagecollected/](http://code.google.com/p/garbagecollected/)

Yesterday evening I came up with an interesting approach for implementing [Josh Bloch's revised GoF Builder pattern](http://developers.sun.com/learning/javaoneonline/2007/pdf/TS-2689.pdf) (warning: PDF). After some late night hacking, I can't help but feel that this is very useful stuff. Take a look at Josh's presentation first, and then take a look at this:

```
package builder;

public class SomeObject {
  private final String mandatory;
  private final int optional1;
  private final char optional2;

  private SomeObject (SomeObjectBuilder builder, String mandatory) {
    this.mandatory = mandatory;
    this.optional1 = builder.optional1();
    this.optional2 = builder.optional2();
  }

  public interface SomeObjectBuilder extends Builder {
    SomeObjectBuilder optional1(int optional1);
    SomeObjectBuilder optional2(char optional2);
    int optional1();
    char optional2();
  }

  public static SomeObjectBuilder builder (final String mandatory) {
    return BuilderFactory.make (SomeObjectBuilder.class,
        new BuilderCallback () {
          public SomeObject call (SomeObjectBuilder builder) throws Exception {
            return new SomeObject(builder, mandatory);
          }
    });
  }

  public String toString() {
    return new StringBuilder()
      .append (getClass().getName())
      .append (String.format ("[optional1=%s, ", optional1))
      .append (String.format ("optional2=%s, ", optional2))
      .append (String.format ("mandatory=%s]", mandatory)).toString();
  }

  public static void main(String[] args) {
    System.out.println(SomeObject.builder("Mandatory!")
        .optional1(35)
        .optional2('A')
        .build()
        .toString()
    );
  }
}
```

Console output: `SomeObject[optional1=35, optional2=A, mandatory=Mandatory!]`

Using a dynamic proxy, the `BuilderFactory` provides the `Builder<T>` implementation for a given interface, so that you don't have to write all that horrible boilerplate code. Often you use a builder when constructors get messy, but Builders with many parameters get messy too. Using this approach you not only save time, you also have the advantage of using a static factory method and having your specific builder as an interface instead of a concrete class. Full source code available upon request; feedback/suggestions/improvements appreciated!

Eat that, setter injection ;-)
