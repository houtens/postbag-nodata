# Avatars

The profile pix are not referenced explicitly in the database but the image
files are named after each user id. For example, user id 462 is me, and my
profile picture is `pix/462.jpg`.

The new database has an avatar string field on the users model. A uuid is
generated for the image name and this is stored on the user model.

## Importing images from absp-database

`cmd/pix/main.go` will download the profile pix from the database and save them
to disk. This script accepts two flags: `-fetch` will download profile pix and
save them to disk; `-seed` will generate uuid names for each image, move the
images to the static folder and save the new filename to the user model.


