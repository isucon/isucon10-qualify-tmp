run:
	cargo run --release
build:
	cargo build --locked --release

rundev:
	# systemfd and cargo-watch is required for autoloading
	# % cargo systemfd cargo-watch
	RUST_LOG=info systemfd --no-pid --socket http::1323 -- cargo watch -x run
