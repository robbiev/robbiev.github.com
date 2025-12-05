-JavaFX: The Return of the Bling
May 9, 2007
My summary:
1. Sun started the [JavaFX](http://java.sun.com/javafx/) initiative, which currently includes JavaFX Mobile and [JavaFX Script](http://openjfx.org) (formerly F3). The latter one is a scripting language (and associated libraries) for UI design and eye candy, which builds on Swing and Java 2D foundations. Some interesting bits:

   - It uses the regular Java JAR deployment techniques: standalone, browser applet, Java Web Start.
   - This is Sun's [WPF](http://en.wikipedia.org/wiki/Windows_Presentation_Foundation)/[XAML](http://en.wikipedia.org/wiki/XAML). They even cloned [XAMLPad](http://en.wikipedia.org/wiki/XAMLPad).
   - It looks to me that it's possible to use the JavaFX Script libraries in Java too, because they're regular Java classes. So no obligations there.
   - JavaFX Script ships plugins for Netbeans and Eclipse.
   - The whole thing will be open sourced (GPL)

2. Early next year Sun will release a consumer JRE, which would be a fast and lightweight version of the JRE, targeted at desktop/applet end users. At the same time they will be revamping the installation experience. There's a related project codenamed "Java Kernel" which will modularize the JRE so that Sun can ship minimal JRE versions more easily. In other words, they will try to fix applets.

Obviously they're trying to catch up with initiatives such as Microsoft's Silverlight and Adobe's Flex (and Apollo). But they still have a long way to go. Coming up next: what's missing?
