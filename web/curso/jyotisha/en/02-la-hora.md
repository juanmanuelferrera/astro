# Module 2 — The hour

**Aim:** to understand that the most fragile piece of data in a chart is the hour, and to be
able to convert a clock time into universal time without going wrong.

## The idea

Module 1 left an uncomfortable consequence: if the horizon moves one degree every four minutes,
then **an error of twenty minutes moves the Ascendant five degrees**. An error of two hours
changes its sign. Almost everything that goes wrong in a chart goes wrong here.

There are three different hours, and it pays not to mix them up.

**Clock time (civil time).** What the clock on the wall said. It is a political agreement: it
depends on the time zone that country was keeping that day.

**Universal time (UT).** Greenwich time. It is what the ephemerides use. This is what you have
to arrive at.

**True local time.** What the Sun would mark at that exact point. It almost never matches the
clock, because time zones are wide bands.

## The four traps

**1. Daylight saving.** The worst of them. It changes by country and by year, and the rules have
been altered many times. Spain has been on Berlin time since 1940. Many countries applied it in
odd years during wars and oil crises. **Never assume it: check it for that country and that
year.**

**2. The time zone is not the longitude.** Almost all of Spain lies west of Greenwich, yet it
keeps central European time. A birth in Galicia carries more than two hours of gap between the
clock and the Sun.

**3. The rounded hour.** "He was born around two" is not a piece of data: it is a half-hour
interval, and half an hour is about seven degrees of Ascendant. You have to ask where the hour
comes from: the hospital record, the civil register, or somebody's memory.

**4. Midnight and midday.** Twelve at night on the 3rd is 00:00 on the 3rd, not on the 4th. It
is a silly mistake and it is made constantly.

## How an hour is validated

Asking "what time were you born?" is not enough. Almost everybody answers with an hour they have
never checked, inherited from a family conversation. Validating it is archive work, and it is
done before calculating anything.

**The sources, best to worst:**

1. **The hospital record.** It logs the minute. It is the most reliable thing there is. In Spain
   anyone can request their own medical record from the hospital where they were born, though
   old archives are sometimes gone.
2. **A literal birth certificate** from the Civil Register. It carries the hour the hospital
   declared. In Spain it is requested free through the Ministry of Justice website. This is the
   practical source: documentary, reachable, and with an hour on it.
3. **The family book.** Sometimes it has the hour, sometimes not.
4. **Something written at the time** — a diary, a letter, a telegram. Worth more than a memory
   because it was written while it was fresh.
5. **Somebody's memory.** The worst source, and the one used ninety per cent of the time.

**Signs that an hour has been invented:**

- **It ends in 00 or 30.** If you ask several people and a lot of round hours come back, those
  hours are estimates. Births do not cluster on the hour and the half hour.
- **It is the registration hour, not the birth hour.** It is written down when somebody fills in
  the form, which may be hours later. Watch out for office hours.
- **Two relatives give different hours.** Ask each one separately and without prompting —
  saying "was it around three?" contaminates the answer.

The other way round, a **scheduled caesarean** usually carries a very reliable hour: it was on
the surgical record.

**And if it cannot be validated**, the right answer is not to pick a pretty hour. It is to work
with an **interval**: "between 9 and 11". Then look at what changes inside that interval — if
nothing relevant changes, go on; if the Ascendant changes sign, it has to be rectified before
you read anything about houses. That is module 4.

And it has to be said out loud to the client: **with a doubtful hour, the planets in signs still
hold and the houses do not.**

## The procedure

1. Clock time and date, as they were recorded
2. Was daylight saving in force that day in that country? If so, subtract an hour
3. Apply the standard zone of the place → you get **UT**
4. Check: does the result fall on the right date? Sometimes the day changes

## Questions

1. A child is born in Barcelona on 15 July 1985 at 03:20. What is the UT?
   *(Summer → UTC+2. UT = 01:20 the same day.)*

2. Someone is born in Madrid at 00:30 on 1 January. Which date do you use for the ephemeris?
   *(UTC+1 in winter → UT = 23:30 on 31 December. The day and the year both change.)*

3. The mother says "it was mid-morning". What Ascendant margin are you handling?
   *(From 09:00 to 12:00 is three hours: about 45°, that is, up to a sign and a half. A chart
   cannot be put up with that. Either rectify or get the data.)*

4. Why is the longitude of the place not enough to know the time zone?
   *(Because the zone is an administrative decision, not an astronomical one.)*

## Exercise

Have them reconstruct the conversion of **their own** hour, step by step, and compare it with
what `carta.py --pasos` prints on lines 1 and 2. If it does not match, have them find which of
the four traps they fell into.

And the archive work, which is done outside the session: **trace their hour up to the highest
source they can reach** on the list above. Have them come back saying which document it comes
from, not who told them. If it ends at "my mother told me", have them ask somebody else
separately and compare.

## Closing question

**Why does a chart with an approximate hour still serve for some things and not for others?**

They should get to: the planets in signs remain valid; the angles and the houses do not. Which
means you can talk about temperament and not about the areas of a life.

## Common mistakes

- Applying today's daylight saving to a birth fifty years ago
- Adding where you had to subtract. West of Greenwich you add to reach UT; east of it you subtract
- Accepting a "family" hour without tracing it to a document
- Suggesting the hour to the relative as you ask, and keeping whatever comes back
- Confusing the hour of birth with the hour it was entered in the register
