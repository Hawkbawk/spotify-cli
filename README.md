# spotify-cli

## About

A (barebones) Spotify CLI, done for Hackweek at Instructure, in Go.

Currently, supports searching for tracks and artists, browsing
and playing the current user's playlists, adjusting playback state
(play, pause, next, previous), and transferring playback to another device.

## Installation

To install from source, a working Go installation and a registered Spotify
application is required. See [the Go website](https://golang.org/) for
instructions on how to download the latest version of Go. To register an
application with Spotify, visit their [website](https://developer.spotify.com/my-applications/).

You'll get an client ID and a secret key for your app. Store this information
in the SPOTIFY_ID and SPOTIFY_SECRET environment variables so that the
application can take actions on your behalf. If you want a method to securely
store and access this info, which you really should, see
[Vaulted](https://github.com/miquella/vaulted).

Once that's done, clone the repo:

`git clone github.com/Hawkbawk/spotify-cli`

Then, move into the directory and run go install:

`cd spotify-cli && go install`

From there, assuming you have your GOBIN environment variable set, simply run
`spotify-cli` from any bash terminal to start the command line application. Use
the arrows keys or VIM movement keys to select an option, "/" to search
(on menu's that support it), and enter to confirm. That's it! You're ready to
control your Spotify experience, all from the command line.

## Note

It's important to note that you have to have a version of the Spotify app
running on one of your devices, as Spotify's Web API doesn't currently support
directly streaming tracks to a registered application, which is fair, but still
no fun for developers who just want to tinker.

## Reporting Bugs & Requesting Features

This is a Hackweek project, so bugs are fully expected. This was written in a
week after all, during an internship no less. However, feel free to submit
any bug reports/feature requests on the project's GitHub page. If you even
want to try fixing the bug you notice, go ahead and submit a pull request!
