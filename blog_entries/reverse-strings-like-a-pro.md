Reverse Strings Like a Pro
March 15, 2013
After quitting my job in February I now started thinking about all the fun I'll
have interviewing for a new one. It's not all bad of course, but I'm not
particularly looking forward to writing yet *another* anagram
detector or string reversal function. But this did make me wonder: how does
string reversal work with [Unicode](http://en.wikipedia.org/wiki/Unicode)? For example, Java's `String` is made up of `chars`.
The `char` data type is 2 bytes. But **not all visible characters fit in 2 bytes**:

- Anything outside of the [Basic Multilingual Plane (BMP)](http://en.wikipedia.org/wiki/Plane_%28Unicode%29)
- [Combining characters](http://en.wikipedia.org/wiki/Combining_character)

The reverse is sometimes also true. For [ligatures](http://en.wikipedia.org/wiki/Typographic_ligature),
multiple visible characters (also called [graphemes](http://en.wikipedia.org/wiki/Grapheme))
can map to a single `char`.

What does all this mean? **Reversing a String is not like reversing an array.**

Take for example the word
&#x0041;&#x0301;&#xFB03;cion&#x1D41A;&#x1E0B;&#x0323;o. This is not exactly the
[original spelling](http://en.wiktionary.org/wiki/afficionado); but
let's see how it looks in terms of Unicode code units
(`chars`) grouped by Unicode code point (individual "building blocks"
of the word):
<table>
  <tr>
    <td>&#x0041;</td>
    <td>&#x0301;</td>
    <td>&#xFB03;</td>
    <td>c</td>
    <td>i</td>
    <td>o</td>
    <td>n</td>
    <td>&#x1D41A;</td>
    <td>&#x1E0B;</td>
    <td>&#x0323;</td>
    <td>o</td>
  </tr>
  <tr>
    <td>0x0041</td>
    <td>0x0301</td>
    <td>0xFB03</td>
    <td>0x0063</td>
    <td>0x0069</td>
    <td>0x006F</td>
    <td>0x006E</td>
    <td>0xD835 0xDC1A</td>
    <td>0x1E0B</td>
    <td>0x0323</td>
    <td>0x006F</td>
  </tr>
</table>

Initially you might have thought that
"&#x0041;&#x0301;&#xFB03;cion&#x1D41A;&#x1E0B;&#x0323;o".length()
would have returned `11`. But it doesn't; it actually returns
`12`. This is because **`String#length()` returns the number
  of `chars`, not the number of visible characters**. Things to note:
- Two visible characters (&#x0041;&#x0301; and &#x1E0B;&#x0323;) have combining characters, thus needing 2 `chars` each
- &#xFB03; is a ligature; it has 3 visible characters, yet only needs 1 `char`
- &#x1D41A; is outside of the BMP, so needs 2 `chars`

Now you can probably imagine the fun Twitter has [enforcing the 140
  character limit](https://developer.twitter.com/en/docs/basics/counting-characters). Anyway, knowing all this, we can do the math for our
example: `2 + 1 + 1 + 1 + 1 + 1 + 2 + 2 + 1 = 12`.

So how about reversing our example? Here are the options I've identified:

1. **Reverse the `char` array**. This is what most people do when they're asked to
  write a string reversal function during an interview. But for all the reasons I
  mentioned above this approach has severe limitations. Though it has to be said
  that if you're not reversing a Unicode string, this is probably the most
  efficient way to do it. But that's not the goal at the moment.
2. Use **[StringBuilder#reverse](http://docs.oracle.com/javase/7/docs/api/java/lang/StringBuilder.html#reverse%28%29)**
  which reverses the `char` array with support for code points outside of the BMP
  (identifying [surrogate
    pairs](http://en.wikipedia.org/wiki/UTF-16#Code_points_U.2B10000_to_U.2B10FFFF), so occurances of more than one `char` to represent a unicode code
  point). This does not reverse combining characters correctly.

3. **Normalize**, then use
  **`StringBuilder#reverse`**. Use the
  [Normalizer](http://docs.oracle.com/javase/7/docs/api/java/text/Normalizer.html) on your String to get rid of as many combining
  characters as possible before reversing. This helps
  but does not work for all character combinations. Also this
  has the effect of modifying your data (e.g. you can't reverse the reverse to
  get back the original).

4. Use **[BreakIterator](http://docs.oracle.com/javase/7/docs/api/java/text/BreakIterator.html)**. This is the best way I have found to reverse a Unicode
  String and also the only way to correctly reverse our example. Note that there is
  still a limitation with ligatures; those letters will obviously not get
  reversed. If this is needed you'll probably want to use the Normalizer after
  all, pehaps on a subset of your data.

Below you can find an example of using `BreakIterator`. This
correctly reverses our example (except the ligature) and is a round-trip
function (meaning you can reverse the reverse and get the original input back).

```
public static String reverseUnicode(String source, Locale locale) {
    BreakIterator boundary = BreakIterator.getCharacterInstance(locale);
    boundary.setText(source);

    char[] reversedChars = new char[source.length()];

    int end = boundary.last();
    for (int start = boundary.previous(), index = 0;
         start != BreakIterator.DONE;
         end = start, start = boundary.previous()) {
      for(int i = start; i < end; i++) {
        reversedChars[index] = source.charAt(i);
        index++;
      }
    }

    return String.valueOf(reversedChars);
  }
```

I can honestly not remember the last time I needed to reverse a
`String`, but there you go.
