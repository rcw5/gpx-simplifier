Why write this code?
=====

I like to bike tour (See http://www.fuelledbycake.co.uk).

I never plan the entire route before leaving home. Who knows what will happen? We might divert to someplace interesting and all that planning goes out of the window.

So we usually take a small laptop with us, loaded with maps, to plot routes which are then transferred to my Garmin GPS (an Edge 800).

This tablet weighs over a kilo, plus a charger, takes up space in a pannier and the battery needs to be kept topped up. It's a pain.

So I am experimenting using an iPad or even my iPhone, planning routes on http://cycle.travel, then saving them to a MicroSD card for the Garmin. I plan to write up how I do this as it took a little trial and error.

cycle.travel lets you plan long distance journeys and export them as a GPX, but doesn't support splitting them into smaller pieces. I want to do this because Garmins don't like large tracks with lots of trackpoints. Too many and it gets confused and can truncate the file, so you lose a portion of the route.

This simple app takes a GPX file and will split it into `n` files, then simplify each file so it has at most `m` points.

It outputs a zip file containing the now-split files ready to upload to your GPS device.
