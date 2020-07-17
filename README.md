# Mattermost jobpost-plugin

_**A Plugin that makes task of creating Jobposts in Scaler chat based on Mattermost easier and cleaner**_

### Usage

* `/jobpost` - opens up an [interactive dialog] to create a Jobpost
* `/remind list` - displays a list of jobposts created by you

### Build

```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.github.rr250.jobpost-plugin-0.0.1.tar.gz
```
