# imgsort

## Command line utility to organize images using EXIF metadata

This program recursively traverses a directory and moves all images to a
fixed folder structure (YYYY/month/DD/file). To do this, I grab the date the
image was created from the EXIF metadata. If this doesn't exist in the image,
I move it to a default date (some time in 2000). Currently I only move files
detected to be images, but I may add the ability to recognize other files.

In my cursory testing, I've found the program capable of fully traversing and
organizing 290 mb of pictures in about 4 seconds.