# Astria Indexer

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