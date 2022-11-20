
KEY_FILE=server.key
CERT_FILE=server.crt
HOST="opendev.dev"
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
# kubectl create secret tls simpsoncrt1 --key keyfile1.key --cert cerfile.crt

openssl x509 -inform der -in server.cer -out certificate.pem
openssl x509 -in server.crt -out server.pem -outform PEM
