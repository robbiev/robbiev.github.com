Defining Quality
March 1, 2013
A while ago a saw [this interesting post](http://antirez.com/news/43) on reliability of software (in this case Redis). The author says:

> Software reliability is an incredibly complex mix of ingredients.
>
> 1. Precision of the specification or documentation itself. If you don't know how the software is supposed to behave, you have things that may be bugs or features. It depends on the point of view.
> 2. Amount of things not working accordingly to the specification, or causing a software crash. Let's call these just software "errors".
> 3. Percentage of errors that happen to be in the subset of the software that is actually used and stressed by most users.
> 4. Probability that the conditions needed for a given error to happen are met.

You can also think of this as a definition of quality. Quality is just something that people perceive. Your users perceive it based on their expectations,  the analyst perceives it based on their specification and you as a programmer perceive it based on your software. This is worth capturing for future reference, so I decided to create the following diagram. What is quality?

![Quality diagram](/img/quality.png)
