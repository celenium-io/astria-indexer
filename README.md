# Astria Indexer
[![Build Status](https://github.com/celenium-io/astria-indexer/workflows/Build/badge.svg)](https://github.com/celenium-io/astria-indexer/actions?query=branch%main+workflow%3A%22Build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Coverage](https://github.com/celenium-io/astria-indexer/wiki/coverage.svg)](https://raw.githack.com/wiki/celenium-io/astria-indexer/coverage.html)

## Standalone

To run indexer clone repository, open directory containing project and run build command. Example:

```bash
git clone git@github.com:celenium-io/astria-indexer.git
cd astria-indexer
make build
```

or by docker compose


```bash
git clone git@github.com:celenium-io/astria-indexer.git
cd astria-indexer
docker-compose up -d --build
```

You have to set environment variables for customizing instances and indexing logic. Example of environment file can be found [here](.env.example).

You can run full stack with GUI:

```bash
export TAG=v1.0.0 # set here needed backend version
export FRONT_ENV=v1.0.0  # set here needed frontend version
docker-compose -f full.docker-compose.yml up -d 
```

After that you can find GUI on `http://localhost:3000`.