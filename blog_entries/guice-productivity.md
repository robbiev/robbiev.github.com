-Guice Productivity
May 8, 2008
Plugin hell: we've all been there. Large frameworks like Spring, Seam, ... often require IDE plugins for you to be productive. That's cool. But installing plugins and getting them to work right, at least when using Eclipse, is utter, utter torture.
Some of them don't work as advertised. Some of them don't work with your particular Eclipse version. Some of them make others crash.

What's refreshing about a framework like Guice is that you don't need all of that. Because Guice is an all-Java framework, you don't *need* to install any IDE plugins. Your Java IDE _is_ the plugin. There's years and years of hard work at your fingertips, so why not work _with_ it instead of working against it? Developing with Guice feels like coming home.

Recently I cooked up [a screen cast](http://users.telenet.be/robbiev/multibind/) that shows off [Multibindings](http://publicobject.com/2008/04/guice-multibinder-api-proposal.html), a much anticipated Guice 2.0 feature. As I mention in the beginning, I use some IDE templates for my Guice development, and you'll probably agree with me that it shows. I'll share those with you in a minute. But first, my brain fart of the day: **Know and learn your IDE, and expect the same from the frameworks you use.** 70% of the code generation refactorings/templates in [my Multibindings screencast](http://users.telenet.be/robbiev/multibind/) are built right into the IDE. And 100% of them are configurable out of the box.

![](/img/gah.png)

Here's what I've currently set up. This configuration is Eclipse specific, but I'm sure IDEA and others have similar features.
- **Static imports:** *Window => Preferences => Java => Editor => Content Assist => Favorites*.
  - This thing is pure gold. It allows you to specify types that have static methods on them, so that Eclipse can index them and suggest static imports for methods you often use. This feature is an awesome time saver for things like Google Collections, but also for JUnit tests and not surprisingly Guice. For Guice I've added the `Matchers` class (AOP) and the `Names` class (eeeviiiil).

- **Templates:** *Window => Preferences => Java => Editor => Templates*. Using this simple IDE feature, all boilerplate is just a `CTRL+SPACE` away. Important: press save (`CTRL+S`) and Organize Imports (`CTRL+SHIFT+O`) after you complete a template. This makes sure that you have all the needed import statements.
  - *Binding Annotations:* creating these is tedious and error prone. One approach is to copy-paste the boilerplate, the other is to have a template. I've set one up that automatically inserts the annotation headers you'd usually want.
  - *`Injector` creation*: I tend to run a lot of "experiments", also known as the "Launching Main893.java" syndrome. So next to the binding annotation header, I've also set up a template that creates me an Injector with an inline `AbstractModule`. Also works great for demos.
  - *Constructor generation:* simple template that creates the default constructor and puts your cursor in between the parenthesis so you can immediately start typing.

![](/img/gah_after.png)

Other notable built-in features are `CTRL+1` (Quick Fix) on constructor fields, which can create private final instance variables for selected constructor arguments, and typing a method name and hitting `CTRL+SPACE` to override or implement that method.

Now without further ado, here are some of my templates.
>
Mapped to 'gah' (as in Guice Annotation Header):

```
@Retention(RetentionPolicy.RUNTIME)
@Target({ElementType.FIELD, ElementType.PARAMETER})
@BindingAnnotation
```

Mapped to 'gin' (as in Guice INjector and booze):

```
Injector i = Guice.createInjector(new AbstractModule() {
    protected void configure() {
        ${cursor}
    }
});
```

Mapped to 'constr' (creates a constructor and sets the cursor where it needs to be):

```
public ${enclosing_type}(${cursor}) {
}
```

Enjoy!
