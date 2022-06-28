# RSS feeds to Wallabag .::|WIP|::.

Take rss feeds and add them to [wallabag](https://www.wallabag.it/en).

## Setup

Copy the example config from this repo and add it to `$HOME/.local/share/go-rss-to-wallabag/config.yaml`

Fill in the following pieces of information :

- client_id
- client_secret
- username
- password
- baseurl

This information can either be a string or a command to obtain that string (e.g. a `pass` command). A command should be supplied using `$(<command>)` in place of the double quotes.
