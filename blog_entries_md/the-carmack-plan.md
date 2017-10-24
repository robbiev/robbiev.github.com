The Carmack Plan
October 24, 2017

A couple of years ago I read *[Masters of Doom](https://www.amazon.com/dp/0812972155)*, which tells the story of [id Software](https://www.idsoftware.com) (of [Doom](https://en.wikipedia.org/wiki/Doom_(series)) and [Quake](https://en.wikipedia.org/wiki/Quake_(series)) fame). Programmer John Carmack had an interesting approach to keeping track of his daily work. From the book:

>As life in the war room pressed on, Carmack took it upon himself to let
gamers know that, yes, id really was moving along with its work on Quake.
So he decided to upload his daily work log, or, as it was known, a .plan file, to
the Internet. Plan files were often used by programmers to keep each other
informed of their efforts but had yet to be exploited as means of communicating
with the masses. But id’s fans had suffered months, years, of Romero’s
unsubstantiated hyperbole, Carmack felt and it was time that they saw some
hard data.

To access Carmack's `.plan` you would use the [finger protocol](https://en.wikipedia.org/wiki/Finger_protocol). Here's what you would have seen when you viewed his `.plan` file, again from *Masters of Doom*:

```
[idsoftware.com]
Login name: johnc In real life: John Carmack
Directory: /raid/nardo/johnc Shell: /bin/csh
Never logged in.
Plan:

This is my daily work ...

When I accomplish something, I write a * line that day.

Whenever a bug / missing feature is mentioned during the day and
I don’t fix it, I make a note of it. Some things get noted many times
before they get fixed.

Occasionally I go back through the old notes and mark with a +
the things I have since fixed.

--- John Carmack

= feb 18 ===================================
* page flip crap
* stretch console
* faster swimming speed
* damage direction protocol
* armor color flash
* gib death
* grenade tweaking
* brightened alias models
* nail gun lag
* dedicated server quit at game end
+ scoreboard
+ optional full size
+ view centering key
+ vid mode 15 crap
+ change ammo box on sbar
+ allow “restart” after a program error
+ respawn blood trail?
+ -1 ammo value on rockets
+ light up characters
```

Many more examples of `.plan` entries can be found [online](https://github.com/ESWAT/john-carmack-plan-archive). Here's one from [August 12 1997](https://github.com/floodyberry/carmack/blob/fc09ed3e7dde67b296a6840524a7c9ec8c36511a/plan_files/johnc_plan_19970812.txt):

```
* qe4 project on command line
* qe4 rshcmd replacement
* qe4 select face
* qe4 avoid multiple autosaves
* qe4 region selected brushes
* bindlist command
* imagelist command in ref_soft

+ leaktest
+ load game.dll from gamedir

pendulum motion
no jump on lava floor?
-game
16 bit wall textures
```

Not all `.plan` entries contain lists of tasks - Some `.plan` entries are [opinion pieces](https://github.com/floodyberry/carmack/blob/fc09ed3e7dde67b296a6840524a7c9ec8c36511a/plan_files/johnc_plan_19990902.txt) and some also include code snippets. I love [the system he used for keeping track of work](https://github.com/floodyberry/carmack/blob/fc09ed3e7dde67b296a6840524a7c9ec8c36511a/page_tools/template/plan.html#L26-L41):

| Prefix    | Meaning                                                                          |
|-----------|----------------------------------------------------------------------------------|
| No prefix | mentioned but not fixed or implemented on that day                               |
| *         | completed on that day                                                            |
| +         | completed on a later day                                                         |
| -         | decided against on a later day                                                   |

What I love about it is the simplicity and that everything in easily `grep`-able plain text. Here are some examples I ran on just the August 12 1997 entry.

Find all open tasks: 

```
$ grep '^[^*+-]' .plan
pendulum motion
no jump on lava floor?
16 bit wall textures
```

Find the last 5 completed tasks:

```
$ grep '^\*' .plan | head -5
* qe4 project on command line
* qe4 rshcmd replacement
* qe4 select face
* qe4 avoid multiple autosaves
* qe4 region selected brushes
```

Alternatively, assuming you keep everything in one file, you can simply open the file in your favourite editor and perform a search there.

Here's what I'm taking away from this:

1. Keep track of things in a way that works for you (or your team).
2. Have your data in a format that is easy to process.

Something to consider before you install enterprise grade work tracking software. John M. Culkin said it best:

>**We shape our tools and, thereafter, our tools shape us.**
