# Gator
Gator is a RSS feed scrapper made following guided project from [bootdev](https://www.boot.dev/lessons/14b7179b-ced3-4141-9fa5-e67dbc3e5242)

## Requirements
For executing this program you need Postgres (ver 15+) installed and Go (1.25+)

## Installing
The isntalling process is fairly simple. Just run 'go install' on the main source code folder.

## Configuration
For configuration you gonna need to create a '.gatorconfig.json' file on your home folder. For linux based programs execute:
`touch ~/.gatorconfig.json`
Add the following data to the file:
`
    {    
        "db_url": "postgres://example"
    }
`
Add your postgres database on the file. Note that sslmode must be disabled for running local.

## Usage
Some of the commands of the system:
`gator login <username>` -> This command log in a user
`gator register <username>` -> Register a user
`gator addfeed <username> <url>` -> Add a feed to a user.
`gator follow <feed_url>` -> Follow a feed for loged in user.
`gator agg` -> Aggregate the feeds that the loged in user follows.