Guice Debug Output
December 8, 2007
Guice has a dirty little secret: it logs timing information to its JDK Logger. A while ago I created a simple utility class for myself to enable or disable the logging of Guice's debug output to the console. Here goes nothing.

```
/**
* Enable or disable Guice debug output
* on the console.
*/
public class GuiceDebug {
    private static final Handler HANDLER;
    static {
        HANDLER = new StreamHandler(System.out, new Formatter() {
            public String format(LogRecord record) {
                return String.format("[Guice %s] %s%n",
                                  record.getLevel().getName(),
                                  record.getMessage());
            }
        });
        HANDLER.setLevel(Level.ALL);
    }

    private GuiceDebug() {}

    public static Logger getLogger() {
        return Logger.getLogger("com.google.inject");
    }

    public static void enable() {
        Logger guiceLogger = getLogger();
        guiceLogger.addHandler(GuiceDebug.HANDLER);
        guiceLogger.setLevel(Level.ALL);
    }

    public static void disable() {
        Logger guiceLogger = getLogger();
        guiceLogger.setLevel(Level.OFF);
        guiceLogger.removeHandler(GuiceDebug.HANDLER);
    }
}
```

Output looks something like:

```
[Guice FINE] Configuration: 51ms
[Guice FINE] Binding creation: 53ms
[Guice FINE] Binding indexing: 0ms
[Guice FINE] Validation: 131ms
[Guice FINE] Static validation: 0ms
[Guice FINE] Static member injection: 2ms
[Guice FINE] Instance injection: 2ms
[Guice FINE] Preloading: 1ms
```

Listen to the Logger, it's making sense. Guice *is* fine! :-)
