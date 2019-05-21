# HAPPY BIRTHDAY
Since I couldn't think of anything more appropriate and suiting, this repository is dedicated to **MAIK** as a birthday present.

## jpeg-filerenamer
Rename JPEG files (`*.jpg` and `*.jpeg`) with their date parsed from EXIF headers. 

The tool 
* ignores files without EXIF headers.
* skips files which already have their proper name.
* checks for potentially existing files with the desired name and behaves in two ways
  * if there is already a file with the given name, the tool compares their MD5 hashes.
    * differing MD5 hashes will cause the tool to suffix the desired name with an automatically increasing number
    * equal MD5 hashes will cause the tool to skip the rename (since a rename to the same name is useless)
  * if there is no file with desired name, renaming will simply take place

The user is asked for confirmation prior to actual processing since all JPEG files in targeted directory will be affected.

#### Usage
* `source` - the source directory to search for JPEG file

#### Why does somebody need this
Assuming you just dropped (and by that also successfully wrecked) an external hard drive with all your precious photos and snapshots of your whole life up until now... **and you didn't have a backup.** 

A data recovery *may* help you getting *some* of the files back, but won't keep you from manually sorting all restored files one-by-one.

In that case, a preceeding basic sorting of photos by their names helps quite a lot.
