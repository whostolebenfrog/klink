#!/bin/bash

# Time for some fun with bash. Use this to ensure that everything is set up
# right, hopefully provide some guidance if not and then finally
# install some dependencies. FUN WITH BASH FOR ALL THE FAMILY.

# make it green
green() {
    echo -ne "\033[032m$1\033[0m"
}

# make it red
red() {
    echo -ne "\033[031m$1\033[0m"
}

# make it gold!
cyan() {
    echo -ne "\033[033m$1\033[0m"
}

# compare two version numbers
vercomp () {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

echo -e "\n\033[033m"
echo "                        888      888 d8b          888"
echo "   _   _                888      888 Y8P          888"
echo "  ( \\_/ )  _   _        888      888              888"
echo " __) _ (__( \\_/ )       888  888 888 888 88888b.  888  888"
echo "(__ (_) __)) _ (__      888 .88P 888 888 888 \"88b 888 .88P"
echo "   ) _ ((__ (_) __)     888888K  888 888 888  888 888888K"
echo "  (_/ \\_)  ) _ (        888 \"88b 888 888 888  888 888 \"88b"
echo "          (_/ \\_)       888  888 888 888 888  888 888  888"
echo "                         ...  ... ... ... ...  ... ...  ..."
echo -e "\n\033[0m"

green "Let's get KLINKING!\n\n"

echo -n "Checking go is on your path: "

command -v go >/dev/null 2>&1 || {
    red "Fail\n\n"
    red "Can't find go on your path. Install go 1.2 then get back to me!\n";
    exit 1;
}

cyan "OK"

echo -ne "\nChecking go version: "

version=`go version | awk '{print $3}' | sed 's/go//'`

echo -ne "Found: $version: "

vercomp $version "1.1.2"

case $? in
        0) cyan "OK\n";;
        1) cyan "OK\n";;
        2) red "Fail\n\n"
           red "Your version of go is too old. Please update to at least 1.1.2 to use klink\n"
           exit 1;;
esac

echo -n "Checking go path: "

if [ -z "$GOPATH" ]; then
    red "Fail\n\n"
    red "Your GOPATH environment variable isn't set.\n" 
    red "Please set your GOPATH, check the readme.md for more information on the format for this.\n"
    exit 1
else
    cyan "OK\n"
fi

echo -n "Checking project is in the correct place: "

klinkpath="$GOPATH/src/nokia.com/klink"
current=`pwd`

if [ $klinkpath != $current ]; then
    red "Fail\n\n"
    red "Your project must be checked out into the src folder, inside your \$GOPATH.\n"
    red "Due to go's package management it must be: \$GOPATH/src/nokia.com/klink\n"
    red "Found: $current\n"
    red "Expected: $klinkpath\n"
    exit 1
else
    cyan "OK\n"
fi

echo -n "Getting dependencies: "

go get "github.com/jmoiron/jsonq"
go get "github.com/jteeuwen/go-pkg-optarg"

cyan "OK\n\n"
green "Finished! You're ready to go\n\n"

exit 0
