# Release Notes
Most recent releases are at the top.

## Version 0.161
 * Version for `bake` command is now optional: will get picked up automatically (with confirmation dialog) from latest Jenkins build as long as:
   - releasePath onix property is setup correctly for the app
   - The Jenkins job that the releasePath corresponds to assigns the version of each release as the build description  
   - An explicit version is not specified to klink
   
