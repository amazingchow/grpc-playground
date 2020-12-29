self-signed:
	openssl req -newkey rsa:4096 -nodes -keyout key.pem -x509 -days 3650 -out cert.pem

clean:
	rm -f cert.pem key.pem

.PHONY: self-signed clean