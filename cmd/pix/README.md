# Avatars

## Description

The absp-database stores profile images of each player. These are not currently
editable in the database and I'm not clear where they came from although it
does look like the London League had something to do with this. We want players to be
able to upload their own photos but we can seed the site with the existing
images if we scrape them.

## Usage

Scraping the avatar images works in two steps - fetch then seed. Fetch will
download all profile pictures and seed with create UUIDs and match them with
users.

## Overview

`-fetch`
- iterate over all users in the postgres database
- x_id is the old user id and references the old image
- try to download all <x_id>.jpg files from the old website
- save to the `seed/data` folder

`-seed`
- iterate again over the database users
- if a download pix exists then
- generate a uuid for the image name
- rename and move the image to avatars
- update the avatars field on the user model

