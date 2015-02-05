# KLINK
Klink is the frontend to our cloud tooling. It makes calls to the various orchestrating applications that allow you to use bake images and deploy applications into aws.

[Release Notes](releases.md)

## What can Klink do?

For a list of commands run `klink` with no arguments.

## Using Klink

To use klink you will first need to have the various other mr-<name> cloud services as klink works by controlling these services. Secondly you will need to alter the conf/conf.go to set the URLs to reflect how you have the various services set up.

At this point you can run klink directly with go run, or compile a static version.

## Developing Klink

**I highly recommend reading through all these instructions before you go any further!**

So you've decided to work on klink? Excellent! There are a few requirements before you get started.

* An installation of go (golang) version `1.1.2` or greater
* go present on your `$PATH`
* your `$GOPATH` environment variable set.
* the source code
* make

### Installing

Go may be available in your package manager of choice e.g. `apt-get install golang` however it is likely that this version is out of date.

I'd suggest getting the latest version from:

http://code.google.com/p/go/downloads/list

an installing with whatever method suits you best. You shouldn't need to compile from source unless you want to locally generate executables for each OS. However there is a job on benkins that can do that for you.

### Go on your path

Add go to your path.

`go version` shows 1.1.2 or greater.

### Set your $GOPATH

OK so your `$GOPATH` variable basically defines your golang workspace. Inside here you will place your source code, dependencies and locally generated executables. I'd suggest something like ~/go

    echo "export GOPATH=$HOME/go" >> .bashrc
    
Or something similar. Don't forget to either source $HOME/.bashrc or start a new terminal to pick that up.

Check with `echo $GOPATH`

### Source

Source code needs to be checked into the **correct package structure** inside your `$GOPATH`

    cd $GOPATH
    mkdir src/nokia.com/
    cd !$
    git clone git@github.brislabs.com:cloud-tooling/klink.git
    
Note that klink requires that it be in src/nokia.com/klink **NOTHING ELSE WILL DO**

That's go for you, it's militant about directory and package structure.

### Make

    cd $GOPATH/src/nokia.com/klink
    make
    
This will ensure that you've done everything above correctly then download the dependencies into your `$GOPATH/src` folder.

#### Fin

Give it a go with:

    go run klink.go

That's it, get hacking. Have fun. Vim highlighting files are inclduded in the golang source if you want them.

## Releasing klink

Run:

`klink build klink`

or

`http://benkins.brislabs.com/job/klink-release/`

or if you have installed go from source and compiled with headers for **ALL** required operating systems you can use the supplied script. Don't do this unless you know what you're doing as it will screw everything up. **MUCH RAGE**.

`./release`
