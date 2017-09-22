# Secure REST Conditional Gateway Recipe
Sample REST conditional gateway with TLS (Transport Layer Security) is enabled along with client authentication.

## Recipe usage instructions

Recipe references below mentioned environment variables which need to be set before creating recipe binary.<br>
SERVER_CERT - Server certificate file <br>
SERVER_KEY - Sever private key file<br>
TRUST_STORE - Directory containing trusted CAs <br>

Now create recipe using mashling cli

```bash
export SERVER_CERT=/etc/ssl/certs/server.crt
export SERVER_KEY=/etc/ssl/certs/server.key
export TRUST_STORE=/etc/ssl/truststore
mashling create -f secure-rest-conditional-gateway.json secureGatewayApp
```

#### openssl can be used if you would like to try with selfsigned certificates

```bash
openssl req \
       -newkey rsa:2048 -nodes -keyout server.key \
       -x509 -days 365 -out server.crt
```

Navigate to gateway bin folder and run the binary. Now the gateway should be running & listening on 3 ports.<br>
9096 -> No security enabled <br>
9097 -> TLS is enabled <br>
9098 -> TLS with client auth is enabled <br>

Use any REST client OR curl to verify the gateway endpoint.

```curl
curl http://localhost:9096/pets/25
```

## If you would like to verify client authentication, use 3rd party go based client - go-mutual-tls

Generate client.crt & client.key for the client, copy them into folder go-mutual-tls/client folder. And update client.go to refer newly generated certificate and key.

Find the line that reads:
```
cert, err := tls.LoadX509KeyPair("../cert.pem", "../key.pem")
```
Change this to
```
cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
```

Copy server certificate into go-mutual-tls/client folder and update client.go

Find the line that reads:
```
clientCACert, err := ioutil.ReadFile("../cert.pem")
```
Change this to
```
clientCACert, err := ioutil.ReadFile("server.crt")
```

Now run the client.

```bash
git clone https://github.com/levigross/go-mutual-tls
cd client
go run client.go
```