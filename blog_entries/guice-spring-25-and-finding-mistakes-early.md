-Guice, Spring 2.5 and Finding Mistakes Early
March 15, 2008
Yesterday I gave a talk on [Guice](https://github.com/google/guice) at the Profict Wintercamp! The goal of the event was to look at how Guice changes the [Spring](https://spring.io/) landscape by letting both sides present their framework. From the Guice side there was me (thanks [Bob](http://crazybob.org), for not being able to come! ;-)), from the Spring side there was [Alef Arendsen](https://twitter.com/alefarendsen).

Surprisingly, there was only one Guice user present! They should have given us an "I went to the Profict Wintercamp and I lived." t-shirt. :-) No seriously, I had a great time and enjoyed talking to other developers and the guys from SpringSource. Let's hope I convinced some people to take a look at Guice!

Right before the talk I had some time to spare, so I quickly threw together a Guice demo application that I then re-implemented using the Spring 2.5 new annotation-driven configuration options. Next, I decided to take a look at how Spring compares with Guice in terms of error detection and error handling. This has always been one of Guice's strengths, but it's also one of those "nice-to-haves". It's like a cell phone. You don't miss it until you have one. All those Spring users don't know what they are missing!

The example code can be found [here](https://github.com/robbiev/garbagecollected/blob/wiki/GuicySpring.md), and the error handling comparison can be found [here](https://github.com/robbiev/garbagecollected/blob/wiki/SpringAndGuiceErrors.md). I'll upload the presentation and the packaged source code soon. Enjoy!
