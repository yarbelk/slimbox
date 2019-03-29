slimbox
=======

Abandoned
---------
I've stopped work on this but don't want to delete it mostly because of the issues; which capture a snapshot of busybox
functionality.  I may come back to this project, but will start the code from scratch as I'm significantlly better
at Golang now; if not better at spelling.

CI
--

[![Build Status](https://travis-ci.org/yarbelk/slimbox.svg?branch=master)](https://travis-ci.org/yarbelk/slimbox)

busybox like project; a paired down version of the gnu tools in on big binary.

Since golang is so good at making tiny binaries, I decided to call it slimbox

Installing
----------

This project is hosted on `gitlab.com`_

.. code::

  go get gitlab.com/yarbelk/slimbox


Contributing
------------

There are many ways to contribute:

Code
~~~~

I have exported all of the busybox helps for the commands (except false and `[` and `[[`.
as Issues.  I will be breaking this down into point release milestones.

Implement one of them, with tests, make a merge request.

Comments
~~~~~~~~

Email me, raise tickets, make documentation pull requests.  Tell me what I'm doing wrong.


implemented
-----------

`cat`, `true`, `false`

note, `cat` is basically the `gnu-tools` version, rather than the busybox.  It does a lot
more because I was experementing with how to do things.
