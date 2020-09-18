all: isuumo

isuumo: app.ts
	deno bundle app.ts > isuumo

run:
	deno run --allow-net --allow-run --allow-read --allow-env --allow-write --unstable ./isuumo

rundev:
	deno run --allow-net --allow-run --allow-read --allow-env --allow-write --unstable ./app.ts
