Remuxing utility helper for mkvmerge
====================================

2does
-----

- [x] Disable verbosity unless -v is defined
- [x] Allow to use without languages setting, appending them all
- [ ] Check files length to ensure all are of the same size, unless param `-S` is specified.
- [ ] Allow to skip duplicated languages, in case there's no -languages specified (the previous setting would add the same audio language for the same file, as different input would probably have the same language multiple times).
- [ ] Be able to specify the proper language id for a track (for cases where language is not properly set in the source)
