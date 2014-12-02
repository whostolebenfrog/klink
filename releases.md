# Release Notes
Most recent releases are at the top.

## Version 0.174
 * Klink ssh now offers a list of servers to pick from
 * klink betabake, bakes against an alternative version of ditto
 * Mac OSX support in autocomplete

## Version 0.163
 * klink ssh <app> <env> now allows you to pick which box to ssh onto rather than just assuming numel 01

## Version 0.161
 * Version for `bake` command is now optional: will get picked up automatically (with confirmation dialog) from latest Jenkins build as long as:
   - releasePath onix property is setup correctly for the app
   - The Jenkins job that the releasePath corresponds to assigns the version of each release as the build description  
   - An explicit version is not specified to klink
