# slimbox
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/yarbelk/slimbox)](https://goreportcard.com/report/gitlab.com/yarbelk/slimbox)

busybox like project; a paired down version of the gnu tools in on big binary.

since golang is so good at making tiny binaries, I decided to call it slimbox

Really tring to do this TDD - and that has resulted in a couple rewrites as I learn better
ways of doing the problem and tdd in its domain.

## Limitations

This is about 5x -20x slower than gnu or busybox.  There are some dubious things
like wc paralellizing out when the bottleneck is in all likelyhood storage iops
not cpu (done non-practical purposes; should be cleaned out).

That said: I intend to improve performance after I get through the first 2 sections
and before I get to runlevel, init and modprobe

## design goals

### libraries

I want to stay stdlib only as much as possible; the major exception is the use of
[pflags](https://github.com/spf13/pflag).  There may be non-stdlib as I get
to more interesting things with network stack or crypto.

### targeted system

I'm just targeting desktop linux right now.

### Layout



```
├── cmd                 // stand alone binary for each command
│   ├── cat
│   │   └── cat_main.go
│   ├── false
│   │   └── false_main.go
│   ├── true
│   │   └── true_main.go
│   └── wc
│       └── wc_main.go
├── lib                 // command logic is in reusable library
│   ├── cat
│   ├── falsy
│   ├── truthy
│   └── wc
├── LICENSE
├── main.go             // fat binary with subcommands main
└── README.md

```

## implemented

See [https://gitlab.com/yarbelk/slimbox/-/boards](kanban board) for where we are.  I want some basic functionality and simple apps and `sh`
such as:

- [x] cat
- [x] wc
- [x] true
- [x] false
- [ ] sh
- [ ] yes
- [ ] cp
- [ ] mv
- [ ] rm
- [ ] ls

Part way through this I intend to implement signal handling as well; once i'm comfortable with how the various programs are structured
then i want some more interesting ones like:

- [ ] df
- [ ] dd
- [ ] ps
- [ ] reboot
- [ ] time  // this is fun because it needs to be consistent
- [ ] mkfifo
- [ ] login
- [ ] sort
- [ ] uniq
- [ ] test

Then the fun ones:

- [ ] exec
- [ ] init
- [ ] runlevel (no idea if i will do this because: its a deap rabit hole)
- [ ] modprobe
