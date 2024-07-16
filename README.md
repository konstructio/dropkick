# dropkick

a utility cli for cleaning up cloud resources

## run it

```
brew tap konstructio/taps
brew install konstructio/taps/dropkick

# run binary to look at what resources will be deleted
dropkick civo --region fra1 

# delete all those resources
dropkick civo --region fra1 --nuke
```
