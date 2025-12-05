-Scripting Guice
May 6, 2007
I really like [Guice](http://code.google.com/p/google-guice/), a Java 5 style Dependency Injection framework by the bright guys at Google.

People who are used to [Spring's DI config](http://static.springframework.org/spring/docs/2.0.x/reference/beans.html#beans-factory-metadata) often complain that they don't like the "compiled configuration" thing (Java Modules) Guice has. In Spring, all your wiring is usually done in XML. So on one side you have compiled Java Guice Modules, a more natural programming model for Java developers, and on the other side you have XML configuration, nicely externalized.  And yes, I know the Spring guys plan on having a Java version of their configuration syntax, but it looks just like the XML, so I don't see any real advantages there.

Anyway, for the Spring users who want to try out Guice, you could externalize your configuration to a properties file, and load it from that. Or whatever format you're willing to use. Heck, you could even write some code that parses a Spring config file and translates it into a Guice Module.

But here's another idea: what about using a scripting language for all your object wiring needs? Why create a configuration file if you could use the highly readable Guice Binder syntax? Well, I tried it, and it  works great :-)

Let's start with an example. This one's about me drinking some beer. So here's me:

```
public class Robbie {
    private Beer myBeer;

    @Inject
    public Robbie(@Strong Beer freshBeer) {
        this.myBeer = freshBeer;
    }
    public void startDrinking() {
        myBeer.drink();
    }
}
```

And here's some beer:

```
public interface Beer {
    void drink();
}
public class Duvel implements Beer {
    public void drink() {
        System.out.println("Duvel!");
    }
}
public class StellaArtois implements Beer {
    public void drink() {
        System.out.println("Stella Artois");
    }
}
```

The `@Strong` annotation basically means I'd like a strong beer, not just some water with taste. So I created a Guice `BindingAnnotation` for that, and expressed the binding in a module:

```
@Retention(RetentionPolicy.RUNTIME)
@Target({ElementType.FIELD, ElementType.PARAMETER})
@BindingAnnotation
public @interface Strong {}
```

```
public class BeerModule implements Module {
    public void configure(Binder binder) {
        binder.bind(Beer.class).to(StellaArtois.class);
        binder.bind(Beer.class).annotatedWith(Strong.class).to(Duvel.class);
        System.out.println("Configured using Java implemented Guice Module!");
    }
}
```

So everyone is getting regular beers, unless they specify they want some stronger beer. You could run this code like this:

```
public class StartGuice {
    public static void main(String[] args) {
        Injector i = Guice.createInjector(Stage.DEVELOPMENT, new BeerModule());
        Robbie robbie = i.getInstance(Robbie.class);
        System.out.print("Robbie starts drinking: ");
        robbie.startDrinking();
    }
}
```

This prints:

```
Configured using Java implemented Guice Module!
Robbie starts drinking: Duvel!
```

Hmm... tasty. Now, let's get us some scripting. I'll thankfully use the new [Java 6 scripting support](https://java.net/projects/scripting/), you could probably use the [BSF](http://jakarta.apache.org/bsf/) as well. Also, I'll use [Jython](http://www.jython.org) for scripting, a Java Python implementation. So I throw in `jython.jar`, `jython-engine.jar` (from the scripting site), `aopalliance.jar` next to the `guice.jar`. By the way, if you want to know how all this scripting support stuff works, check out [Jurgen's excellent introduction](http://jroller.com/page/mom?entry=tutorial_running_ruby_using_the).

Anyway, let's try binding our dependencies in Python. Create a file called `BeerBinder.py` that looks like this:

```
import java
from scriptguice import *

    def configure(binder):
        binder.bind(Beer).to(StellaArtois);
        binder.bind(Beer).annotatedWith(Strong).to(Duvel);
        print "Configured Guice using Python method!"
```

Then we add some scripting magic to our Guice Module, so that we delegate the configure call to the Python script:

```
public class BeerModule implements Module {
    public void configure(Binder binder) {
        ScriptEngineManager mgr = new ScriptEngineManager();
        ScriptEngine python = mgr.getEngineByName("python");
        Reader reader = ReadUtil.getReaderForClassPathResource("scriptguice/pythonmethod/BeerBinder.py");
        try {
            python.eval(reader);
            Invocable invocablePython = (Invocable) python;
            invocablePython.invokeFunction("configure", binder);
        } catch (ScriptException e) {
            throw new RuntimeException(e);
        } catch (NoSuchMethodException e) {
            throw new RuntimeException(e);
        }
    }
}
```

If we run `StartGuice` again, the output now looks like:

```
Configured Guice using Python method!
Robbie starts drinking: Duvel!
```

Great, it works! Now let's take it one step further and get rid of the ugly Java / Python mix. Python all the way! Create a file called `BeerModule.py`:

```
import java
from scriptguice import *
from com.google.inject import *

class BeerModule(Module):
    def configure(self, binder):
        binder.bind(Beer).to(StellaArtois);
        binder.bind(Beer).annotatedWith(Strong).to(Duvel);
        print "Configured using Python implemented Guice Module!"

# Factory method that returns new BeerModule
def getBeerModule():
    return BeerModule()
```

So we created a class called `BeerModule`, and subclassed the Java class `com.google.inject.Module`. Jython really starts to shine here. Let's do some magic to get the thing back in Java:

```
public class StartGuice {
    public static Module getPythonBeerModule() {
        ScriptEngineManager mgr = new ScriptEngineManager();
        ScriptEngine python = mgr.getEngineByName("python");
        try {
            Reader reader = ReadUtil.getReaderForClassPathResource("scriptguice/pythonclass/BeerModule.py");
            python.eval(reader);
            Invocable invocablePython = (Invocable)python;
            return (Module)invocablePython.invokeFunction("getBeerModule");
        } catch (ScriptException e) {
            throw new RuntimeException(e);
        } catch (NoSuchMethodException e) {
            throw new RuntimeException(e);
        }
    }

    public static void main(String[] args) {
        Injector i = Guice.createInjector(Stage.DEVELOPMENT, getPythonBeerModule());
        Robbie robbie = i.getInstance(Robbie.class);
        System.out.print("Robbie starts drinking: ");
        robbie.startDrinking();
    }
}
```

It's been a long time since my last beer, so here we go:

```
Configured using Python implemented Guice Module!
Robbie starts drinking: Duvel!
```

Man, how cool is this? Now we have defined the entire Guice module in Python code, and I got drunk in the process. Life just doesn't get any better, does it ;-)

Now you could take it even further and return an array of Modules or what not, which is probably what you want in a real world app. Oh and one final note: don't just copy paste this code. Close the Reader and stuff. I'm just being lazy for the sake of the example.

Thanks to [Bob Lee](https://twitter.com/crazybob) and [Kevin Bourrillion](http://smallwig.blogspot.com/) for creating Guice, the coolest DI framework on the planet.
