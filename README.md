# Astria Indexer
[![Build Status](https://github.com/celenium-io/astria-indexer/workflows/Build/badge.svg)](https://github.com/celenium-io/astria-indexer/actions?query=branch%3Amaster+workflow%3A%22Build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

You can set environment variables for customizing instances and indexing logic. Example of environment file can be found [here](.env.example).
