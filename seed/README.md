# Seeding the database

Handle exporting the old data from mysql and fix numerous issue with the data.


## Schema

The old schema has no FK relations and all joins are done with subqueries! New
schema diagram can be found on dbdiagram
[absp schema](https://dbdiagram.io/d/ABSP_with_relations-61fc255185022f4ee5360682)

## Ratings

Such a mess as we previously included so many varieties of paritial results.
The `type` parameter at least indicates something about the origin of the
records but is overloaded.

## Types

- 0: one or more player is missing (bye), no score or no spread. Action: do not import.
- 1: complete results entered after 2015. Action: import.
- 2: player and spread exist, no scores, recorded as win/loss. Action: import but log scores with one player scoring zero.
- 3: players and score representing games! Recorded as 1-0 or 0-1 simply for win/loss. Action: Record with non-zero spread and aggregate games.
- 4: one or more players is missing. Walkover of non-standard spread but games are not rated. Action: drop.

## End Ratings

The new ratings model will forego the 150 game history and will base new
ratings on predicted and actual wins. We will be adopting the Australian model documented [here](https://scrabble.org.au/ratings/).

The patchy historical game history makes some calculations tricky and we would
like to be able to chart a history of ratings changes per player. The end
ratings calculation makes a reasonable attempt to back-fill historical end
ratings where they have not been recorded.

