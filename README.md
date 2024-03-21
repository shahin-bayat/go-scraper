How to Start?

1. Create a Postgres db

```bash
	 docker run --name scraperdb \
	          -e POSTGRES_PASSWORD=goscraper \
	          -e POSTGRES_USER=shahin \
	          -e POSTGRES_DB=scraperdb \
	          -p 5432:5432 \
	          -d postgres
```

2. Run the server

```bash
  make run
```
Entity relationship diagrams (ERD):

![image](https://github.com/shahin-bayat/go-scraper/assets/51708006/7051dc61-7b4f-4880-9429-efada2ac7d8e)


