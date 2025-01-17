CERTS_DIR = certs
CA_KEY = $(CERTS_DIR)/ca.key
CA_CERT = $(CERTS_DIR)/ca.crt
SERVER_KEY = $(CERTS_DIR)/server.key
SERVER_CSR = $(CERTS_DIR)/server.csr
SERVER_CERT = $(CERTS_DIR)/server.crt

all: $(SERVER_CERT)

$(CERTS_DIR):
	mkdir -p $(CERTS_DIR)

$(CA_KEY): $(CERTS_DIR)
	openssl genrsa -out $(CA_KEY) 2048

$(CA_CERT): $(CA_KEY)
	openssl req -x509 -new -key $(CA_KEY) -sha256 -days 1024 -out $(CA_CERT) -subj "/CN=My CA"

$(SERVER_KEY): $(CERTS_DIR)
	openssl genrsa -out $(SERVER_KEY) 2048

$(SERVER_CSR): $(SERVER_KEY)
	openssl req -new -key $(SERVER_KEY) -out $(SERVER_CSR) -subj "/CN=localhost"

$(SERVER_CERT): $(SERVER_CSR) $(CA_CERT) $(CA_KEY)
	openssl x509 -req -in $(SERVER_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SERVER_CERT) -days 500 -sha256

clean:
	rm -rf $(CERTS_DIR)