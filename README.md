# Snippet DJesus

This project is an implementation based on the guidelines provided in the book "Let's Go" by Alex Eduards. It includes a landing page that displays the existing snippets created on the "Create Snippet" page. To generate a snippet, you must create an account and sign in.

Setup certificates for localhost
```sh
cd tls/ && go run /usr/local/Cellar/go/1.21.5/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Create `.envrc` file based on the existing `.envrc.example`.
```sh
cp .envrc.example .envrc
```

Run DB migrations
```sh
make run migrations/up
```

Finally, to run the app you simply need to run the following command.
```sh
make run
```
