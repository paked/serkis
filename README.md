# Serkis

**this README.md is a work in progress**

Serkis is a simple and powerful personal wiki. It's built on `git` and `markdown`, and wants to fit in with *your* notetaking process without imposing any of it's own opinions.

## Where does Serkis fit into my workflow?

Serkis assumes that you already do something like the following:

1. Have a git repository of all your personal notes
2. Edit them locally

It picks up from where you left off, allowing you to view and edit old notes, and create new ones online -- with absolutely everything being driven off of your git repository.

## Configuration

Here's an example configuration for you to check out:

```json
{
  "http-username": "fiona",
  "http-password": "confucianism",
  "git-url": "https://github.com/fionahackworth/brain.git",
  "git-username": "fionahackworth",
  "git-password": "octet",
  "git-author-name": "Fiona Hackworth",
  "git-author-email": "fiona@hackworth.com",
  "github-webhook-secret": "something-here"
}
```

## Installation

The recommended way to deploy Serkis is to use `docker` and `docker-compose`. By default it serves on port `8765`:

```
cd serkis
docker-compose build run # run is the compose service with the production ready config in it
docker-compose run --service-ports -d run # make sure docker binds to ports, and use -d so that the dontainer runs in the background
```

Congratulations! If you configured everything correctly, your personal notes will now be live on Serkis. [*You probably now want to put it behind Nginx or another reverse proxy to handle HTTPS.*](https://github.com/paked/schmidt/blob/master/sites-available/ascii.schmidt.harrison.tech)

## Things it does (or will do):

1. [x] Work interoperably with your normal Markdown notetaking habits
2. [x] Put your notes behind an auth wall
3. [x] Allow you to edit your files online (via mobile, or desktop, or potato-with-a-brower-and-internet-connection)
4. [ ] Share specific files with your friends or co-workers or parents
