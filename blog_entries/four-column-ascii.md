Four Column ASCII
January 31, 2017

I found [this gem](https://news.ycombinator.com/item?id=13499386) on [Hacker News](https://news.ycombinator.com/item?id=13498365) the other day. User [soneil](https://news.ycombinator.com/user?id=soneil) posted to a [four column version](http://pastebin.com/cdaga5i1) of the ASCII table that blew my mind. I just wanted to repost this here so it is easier to discover.

Here's an excerpt from the comment:

>I always thought it was a shame the ascii table is rarely shown in columns (or rows) of 32, as it makes a lot of this quite obvious. eg, http://pastebin.com/cdaga5i1
It becomes immediately obvious why, eg, ^[ becomes escape. Or that the alphabet is just 40h + the ordinal position of the letter (or 60h for lower-case). Or that we shift between upper & lower-case with a single bit.

You know in [ASCII](http://www.asciitable.com/) there are 32 characters at the beginning of the table that don't represent a written symbol. Backspace, newline, escape - that sort of thing. These are called [control characters](https://en.wikipedia.org/wiki/Control_character).

In the terminal you can type these control characters by holding the `CTRL` (control characters, get it?) key in combination with another key. For example, as many experienced vim users know pressing `CTRL+[` in the terminal (which is `^[` in [caret notation](https://en.wikipedia.org/wiki/Caret_notation)) is the same as pressing the `ESC` key. **But why is the escape key triggered by the `[` character? Why not another character?** This is the insight soneil shares with us.

Remember that ASCII is a 7 bit encoding. Let's say the following:

* The first two bits denote the group of the character (2^2 so 4 possible values)
* The remaining five bits describe a character (2^5 so 32 possible values)

In the linked table, which I reproduce below, the four groups are represented by the columns and the rows represent the values.

|00 |01 |10 |11 |   |
|---|---|---|---|---|
NUL|Spc|@|\`|00000|
SOH|!|A|a|   00001|
STX|"|B|b|   00010|
ETX|#|C|c|   00011|
EOT|$|D|d|   00100|
ENQ|%|E|e|   00101|
ACK|&|F|f|   00110|
BEL|'|G|g|   00111|
BS |(|H|h|   01000|
TAB|)|I|i|   01001|
LF|*|J|j|    01010|
VT|+|K|k|    01011|
FF|,|L|l|    01100|
CR|-|M|m|    01101|
SO|.|N|n|    01110|
SI|/|O|o|    01111|
DLE|0|P|p|   10000|
DC1|1|Q|q|   10001|
DC2|2|R|r|   10010|
DC3|3|S|s|   10011|
DC4|4|T|t|   10100|
NAK|5|U|u|   10101|
SYN|6|V|v|   10110|
ETB|7|W|w|   10111|
CAN|8|X|x|   11000|
EM |9|Y|y|   11001|
SUB|:|Z|z|   11010|
|**ESC**|;|**[**|{|**11011**|
FS|<|&bsol;|\||  11100|
GS|=|]|}|    11101|
RS|>|^|~|    11110|
US|?|_|DEL|  11111|

Now in this table, look for ESC. It's in the first group, fifth from the bottom. It's in the first column so its group has bits '00', the row has bits '11011'. Now look on the same line, what else is there? Yep, the '[' character is there, be it in a different column:

* `10 11011` means [
* `00 11011` means ESC

So when we you type `CTRL+[` for `ESC`, you're asking for the equivalent of the character `11011` (`[`) out of the control set. Pressing CTRL simply sets all bits but the last 5 to zero in the character that you typed. You can imagine it as a bitwise AND.

```
  10 11011 ([)
& 00 11111 (CTRL)
= 00 11011 (ESC)
```

This is why `^J` types a newline, `^H` types a backspace and `^I` types a tab. This is why if you `cat -A` a Windows text file, it has ^M printed all over (meaning `CR`, because newlines are `CR+LF` on Windows).
