# casio-f91w-go

Emulates hourly notification beep of the classic Casio F91W wrist watch.

# Install

    $ go install github.com/MawKKe/casio-f91w-go@latest

# Usage

Normal usage:

    $ casio-f91w-go 

this will launch the program in continuous mode which plays the "beep beep" notification sound every hour.

You can debug playback issues with

    $ casio-f91w-go -debug

which will bring the beeper interval down to ten seconds.

# License

Copyright 2022 Markus Holmström (MawKKe)

The works under this repository are licenced under Apache License 2.0. See file LICENSE for more information.


# Contributing

This project is hosted at https://github.com/MawKKe/casio-f91w-go

You are welcome to leave bug reports, fixes and feature requests. Thanks