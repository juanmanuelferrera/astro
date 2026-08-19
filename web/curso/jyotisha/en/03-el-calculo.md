# Module 3 — The calculation

**Aim:** to put up a chart **by hand**, once in a lifetime. Afterwards they will use software,
but knowing what the software is doing.

This module is slow. Do not speed it up. It is the one that separates someone who knows
astrology from someone who knows how to work a program.

## The idea

Calculating a chart means answering a geometrical question: **which degree of the ecliptic were
the horizon and the meridian passing through at that instant and at that point of the Earth?**

Everything else is tooling to answer it.

## The steps

**1. From civil time to UT.** Already covered in module 2.

**2. Planetary positions.** The ephemeris gives the planets at noon (or at midnight) of each
day. A birth almost never falls there, so you have to **interpolate**: if the Sun advances 57′
that day and the birth was a third of the way through it, it advanced a third of 57′.

This is where the **proportional logarithms** in the old ephemerides come in: they turn the rule
of three into an addition. Explain the mechanism even though it is not needed today.

**3. Sidereal time.** This is the heart of it, and the point where everybody gets stuck.

The solar day lasts 24 h; the sidereal day about 4 minutes less, because while the Earth turns
on itself it also advances along its orbit. **Sidereal time** says which degree of the zodiac is
culminating right now at Greenwich.

You take it from the ephemeris for that day and correct it:
- for the hours elapsed since the tabulated moment
- for the acceleration (about 10 seconds per hour elapsed)
- for the **longitude** of the place: every 15° of longitude is one hour

The result is **local sidereal time**.

**4. Ascendant and Midheaven.** With local sidereal time and **latitude** you enter the Tables of
Houses and read off the two angles. Latitude is needed because the tilt of the horizon against
the ecliptic changes with it: that is why at high latitudes some signs take an age to rise and
others rise in minutes.

**5. Intermediate cusps.** Depending on the system. That is module 6.

**6. Place the planets** in the resulting houses.

## How you correct them

Have them do **every** calculation themselves. When they give a result, run:

```
python3 carta.py --fecha ... --hora ... --tz ... --lat ... --lon ... --pasos
```

That prints UT, Julian day, Greenwich sidereal time, the longitude correction, local sidereal
time, Ascendant and MC.

**Tell them which step they went off at, and by how much. Do not give them the right number.**
Have them redo it from there.

## Questions

1. Why is the sidereal day shorter than the solar day?
2. Two people are born at the same UT, one in Quito and one in Reykjavík. Why is the Ascendant
   different if Greenwich sidereal time is the same?
   *(Because of latitude and longitude: each enters the table on a different row and with a
   different LST.)*
3. If you are 4 minutes out in sidereal time, how far out is the Ascendant?
   *(About one degree, depending on latitude.)*
4. What is latitude for, and what is longitude for?
   *(Longitude: to correct sidereal time. Latitude: to enter the tables of houses.)*

## Exercise

Put up their whole chart by hand. In pencil. Estimated time: one to two hours the first time.
Have them note every step so it can be compared.

## Closing question

**Which piece of data do you need for local sidereal time, and which for the Ascendant?**

## Common mistakes

- Forgetting the acceleration correction
- Adding the longitude where it had to be subtracted
- Entering the tables with the latitude on the wrong side of the equator
- Interpolating the Moon as if it were linear across a whole day: it moves 12-15° daily and
  deserves care
