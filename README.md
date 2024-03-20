# dropkick

a utility cli for cleaning up cloud resources

## run it

```
go build -o dropkick .

# run binary to look at what resources will be deleted
./dropkick civo --region fra1 

# delete all those resources
./dropkick civo --region fra1 --nuke
```