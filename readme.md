Remuxing utility helper for mkvmerge
====================================

Binary to easily generate a Matroska file from multiple inputs.

This ain't a golang library for matroska files, nor an alternative to mkvmerge, but an utility to easily create mkvmerge commands from various input files.

Example
-------

A command like:

~~~bash
remuxing -output output.mkv -languages spa,eng input1.mkv input2.mkv
~~~

Generates this mkvmerge command (expanded here for ease of reading):

~~~bash
mkvmerge \
  -o output.mkv \
  --title  \
  -A -T -S -d 0 input2.mkv \
  -T --default-track 1 --language 1:spa -a 1 --track-name 1: -D -S input1.mkv \
  -T --language 1:eng -a 1 --track-name 1: -D -S input2.mkv \
  -T -s 3 --track-name 3: -D -A --forced-track 3:true input1.mkv \
  -T -s 4 --track-name 4: -D -A input1.mkv \
  -T -s 5 --track-name 5: -D -A input1.mkv
~~~

Command syntax
--------------

~~~bash
remuxing [options] [inputs]
~~~

Note that you can define as many inputs as you want. The input order is important, as it designates files' priority, used to decide between inputs in case both seem to be of the same quality & codec.

### Arguments

- `-v`: Enables verbosity. Optional.
- `-output`: Sets output file. Mandatory.
- `-languages`: Defines the desired output languages. Order is important, first language will be set as default one. Not setting this option will merge all inputs. Optional.

Installing
----------

If you have golang in your system, simply do:

~~~bash
go get github.com/elboletaire/remuxing
go run github.com/elboletaire/remuxing [options] [inputs]
~~~

Otherwise, you can download any of the available built binaries from [the gitlab copy][binaries]:

- Linux [x32][linux x32]/[x64][linux x64]
- Windows [x32][win x32]/[x64][win x64]
- [Mac (x64 only)][mac]

2does
-----

- [x] Disable verbosity unless -v is defined
- [x] Allow to use without languages setting, appending them all
- [x] Add builds for download (using gitlab-ci or drone or...)
- [ ] Check files length to ensure all are of the same size, unless param `-S` is specified.
- [ ] Allow to skip duplicated languages, in case there's no -languages specified (the previous setting would add the same audio language for the same file, as different input would probably have the same language multiple times).
- [ ] Be able to specify the proper language id for a track (for cases where language is not properly set in the source)

[golang]: https://golang.org/
[binaries]: https://gitlab.com/elboletaire/remuxing

[linux x64]: https://gitlab.com/elboletaire/remuxing/-/jobs/artifacts/master/download?job=build%3Alinux-x64
[linux x32]: https://gitlab.com/elboletaire/remuxing/-/jobs/artifacts/master/download?job=build%3Alinux-x32
[win x32]: https://gitlab.com/elboletaire/remuxing/-/jobs/artifacts/master/download?job=build%3Awindows-x32
[win x64]: https://gitlab.com/elboletaire/remuxing/-/jobs/artifacts/master/download?job=build%3Awindows-x64
[mac]: https://gitlab.com/elboletaire/remuxing/-/jobs/artifacts/master/download?job=build%3Amac-x64
