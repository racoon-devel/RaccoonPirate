# RaccoonPirate

![Logo](contrib/logo.jpg)

Application for consuming media content from torrent trackers on-the-fly.

## How it works?

### Content Sources

The application consumes torrent-files or magnet links to peering networks for downloading sought-for content. They are can be provided by one of the following ways:

* discover via [API](https://github.com/RacoonMediaServer/rms-media-discovery/blob/master/api/discovery.yml) by [Racoon Media Server remote backend service](https://github.com/RacoonMediaServer/rms-remote);
* manually add via user interface.

### Delivery Method

The application mounts a directory as Fuse file system and manages I/O operations for the directory. After torrent registration - content directory layout maps to the cache directory. User can use any media player for play the content. Data will be downloading on-the-fly. Data chunk rotation is supported by set data storage limit in config file. 

### User Interface

There are a few frontends enabled:

* cli (coming soon);
* Web UI;
* Telegram Bot integration (coming soon).

## Dependencies

* libfuse2

## Supported Platforms

* Linux;
* Batocera (retro-gaming Linux distro, PirateRacoon supports integration).

## Use Case Scenarios

* Watching any discoverable content on the PC/MiniPC by any media player;
* Install the application on a home server and stream data over network for consuming from various devices (e.g. DLNA servers).  
