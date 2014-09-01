Extremely Functional Java
August 31, 2014
One definition of **functional programming is to write programs that use [a lot of verbs (functions) to manipulate a small amount of nouns (data structures)](http://steve-yegge.blogspot.co.uk/2006/03/execution-in-kingdom-of-nouns.html)**. Here I will explore the following question: **can we write Java in this style**?

First you'd get rid of classic [OOP](https://en.wikipedia.org/wiki/Object-oriented_programming) "objects". You know, the ones in which you mix state and functionality:

```java
class LightSwitch {
  private boolean on;

  Light(boolean on) {
    this.on = on;
  }

  void toggle() {
    this.on = !this.on;
  }

  boolean isOn() {
    return this.on;
  }
}

```

An **easy way** to use a more functional style would be to use a **static method instead of an instance method**. Using static methods is an easy way to decouple the "verb" (the method) from the noun (the boolean). A static method will never be able to access its class' (non-static) fields without receiving it as an argument, as it is associated with the class, not any instance. As an extra we can decide to make the data structure immutable - this will make it thread safe and will eliminate a train of thought when debugging the program (what if the switch was toggled to on here, turned off here...).

```java
class LightSwitch {
  private final boolean on;

  Light(boolean on) {
    this.on = on;
  }
  
  boolean isOn() {
    return this.on;
  }
}

class LightOps {
  static LightSwitch toggle(LightSwitch switch) {
    return switch.isOn() ? new LightSwitch(false) : new LightSwitch(true);
  }
}

```

The problem with this approach, however, is that **your program gets less flexible as it grows**. For example, let's imagine our light switch becomes audio-activated. When I clap my hands I want to toggle the light switch:

```java
class ClapDetector {
  static boolean isClapping(Audio audio) {
    // something crazy
    ...
  }
}

class LightOps {
  static LightSwitch toggleIfClapping(LightSwitch switch, Audio audio) {
    if (ClapDetector.isClapping(audio)) {
      return switch.isOn() ? new LightSwitch(false) : new LightSwitch(true);
    }
    return switch;
  }
}

```

Now what if we wanted to try different audio processing algorithms? Or test the toggling code in isolation? We can't. **The dependency on the `ClapDetector` has been hard-coded**. The only way to really solve this problem pre Java 8 is to switch back to using class instances and interfaces. It is possible to write in a functional style this way, but it is much harder; the temptation to add some "state" here and there will be great. Every function (methods in this case) would also have invisible parameters from the get-go: the instance fields of the surrounding class (dependencies, like `ClapDetector`).

Noticed how I sneakily said *pre Java 8*. After using it for a month or so, I realised that the new **[default methods](http://docs.oracle.com/javase/tutorial/java/IandI/defaultmethods.html) feature is transforming interfaces into ideal holders of functions**. They still can't hold instance data, but now they can hold method implementations. This pretty much transforms interface methods into functions. Also the combination of interface multiple inheritance and default methods turns out to be pretty powerful. Back to the example:

```java
interface ClapDetector {
  boolean ClapDetector$isClapping(Audio audio);
}

interface NaiveClapDetector {
  default boolean ClapDetector$isClapping(Audio audio) {
    // something crazy
    ...
  }
}

interface LightOps extends ClapDetector {
  default LightSwitch toggleIfClapping(LightSwitch switch, Audio audio) {
    if (ClapDetector$isClapping(audio)) {
      return switch.isOn() ? new LightSwitch(false) : new LightSwitch(true);
    }
    return switch;
  }
}

class LightOpsImpl extends LightOps, NaiveClapDetector {}
```

Notice the **dependency injection using standard language constructs**. All I had to do was to create a class in the end which wires up my dependencies the way I want. This actually reminds me of this [excellent presentation](http://www.infoq.com/presentations/post-functional-scala-clojure-haskell) by [Daniel Spiewak](https://twitter.com/djspiewak) in which he details how [Scala](http://scala-lang.org/) traits are really "modules". Java 8 interfaces are indeed very similar to Scala traits.

For those still interested, I explored a slightly larger example of the ideas I described here: https://github.com/robbiev/mars-functional-java

Of course as to be expected, this style also has its flaws. Here are some that I currently see:

* All functions end up in the same namespace. For this reason I prefixed all functions with the interface name in my examples (`InterfaceName$methodName`)
* You can't easily have several different implementations of the same interface in an object graph
* Creating value objects / struct-like data containers is still overly verbose. Perhaps something like [AutoValue](https://github.com/google/auto/tree/master/value) could help here.
