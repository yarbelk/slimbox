# sh

Posix shell implementation

get `nex`
http://www-cs-students.stanford.edu/~blynn//nex/

`go get github.com/blynn/nex`


get the golang `goyacc` tool


```
nex -s Sh sh.nex
goyacc -p "Sh" sh.y
```

I need to setup the generate directives and things to automate this.
