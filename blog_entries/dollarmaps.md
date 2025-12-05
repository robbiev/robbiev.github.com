DollarMaps: Easy Map Creation for Java
April 6, 2008
Inspired by [GQuery](http://timepedia.blogspot.com/2008/04/gwt-road-to-15-language-features-and.html)'s syntax,  I hacked up some code that uses a similar, dollar sign based syntax for creating Java Map instances. And I have to say that it flows really nicely. But first, let's "respect the classics" .

```
Map<Integer, String> map = new HashMap<Integer, String>();
map.put(1, "one");
map.put(2, "two");
map.put(3, "three");

// iterate over all the entries
for (Entry<Integer, String> e : map.entrySet())
    System.out.println(e.getKey() + " " + e.getValue());
```

So, there are three things to notice here.
 - The repeated type parameters on the first line
 - Having to use three lines of extra code to get the items in
 - Having to use getKey() and getValue() in the iteration

If you didn't know yet, most of these issues can be solved by using [Google Collections](http://code.google.com/p/google-collections/). But besides that, I think the dollar sign based syntax, which I called DollarMaps, is slightly more elegant than the usual "let's-do-some-type-inference" factory method.

Creating a HashMap:

```
Map<Integer, String> map = $(1,"blah1")
                          .$(2, "blah2")
                          .$(3, "blah3").asHashMap();
```

Iteration:

```
for(Entry<Integer, String> e : $(1,"blah1").$(2, "blah2")) {
    System.out.println(e.getKey() + " " + e.getValue());
}
```

You can even take iteration further if both the key and the value have the same type. I tried some things out and came up with the double dollar sign syntax to enforce that type rule. So easier, array-based iteration:

```
for(String[] s : $$("1","blah1").$("2", "blah2").asEasy()) {
    System.out.println(s[0] + " " + s[1]);
}
```

I did not run any performance tests (I am creating a 2-element array for each entry).
Anyway, I'll have this up for grabs at my [Google Code project](http://code.google.com/p/garbagecollected/). For the impatient, here's the [download link](https://github.com/robbiev/garbagecollected/raw/master/DollarMaps/dist/DollarMaps-snapshot.zip). Feedback appreciated!
