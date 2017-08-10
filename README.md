PerfectSSLScore
===============

A webserver written in golang to get a perfect ssl test score on https://www.ssllabs.com/ssltest/.

Perfect here means a A+ rating with all subtests beeing 100%.
In order to reach 100% in the 'key-exchange' subtests you must have a RSA4096 bit cert/key pair.
Otherwise you only get 90% in the key-exchange subtests.

Getting a certificate with lego
-------------------------------

The easiest way to get a certificate is useing [lego](https://github.com/xenolf/lego).

lego -m <your_email> -d <your_domain> -k RSA4096 run

```
lego -d manitu.scusi.io -m flw@posteo.de -k rsa4096 renew
2017/08/10 16:16:44 [INFO][manitu.scusi.io] acme: Trying renewal with 1459 hours remaining
2017/08/10 16:16:44 [INFO][manitu.scusi.io] acme: Obtaining bundled SAN certificate
2017/08/10 16:16:45 [INFO][manitu.scusi.io] acme: Could not find solver for: dns-01
2017/08/10 16:16:45 [INFO][manitu.scusi.io] acme: Trying to solve TLS-SNI-01
2017/08/10 16:16:47 [INFO][manitu.scusi.io] The server validated our request
2017/08/10 16:16:47 [INFO][manitu.scusi.io] acme: Validations succeeded; requesting certificates
2017/08/10 16:16:50 [INFO] acme: Requesting issuer cert from https://acme-v01.api.letsencrypt.org/acme/issuer-cert
2017/08/10 16:16:50 [INFO][manitu.scusi.io] Server responded with a certificate.
```

Clone the git repository
------------------------

```
go get github.com/scusi/PerfectSSLScore
cd $GOPATH/src/github.com/scusi/PerfectSSLScore
```

Build the server
----------------

```
go build -i -v -ldflags="-s -w -X main.version=$(git describe --always --long) -X 'main.buildtime=$(date -u '+%Y-%m-%d %H:%M:%S')'" ./
```


Starting the server
-------------------

```
./PerfectSSLScore -cert .lego/certificates/manitu.scusi.io.crt -key .lego/certificates/manitu.scusi.io.key
```
