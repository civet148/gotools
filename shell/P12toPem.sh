#!/bin/sh

### 第一个参数是要转换的p12文件名，第二个参数是p12文件解密密码 ###
if [ ! $# == 2 ]; then
  echo "Usage: $0 xxx.p12 secret"
  exit
fi

P12CERT=$1
P12PASS=$2
TMP_PASS=123456
PEM_CERT_NAME=${P12CERT%.*}.pem
PEM_KEY_ENC_NAME=${P12CERT%.*}_key.pem
PEM_KEY_NAME=${P12CERT%.*}.key

echo "input cert [$P12CERT] secret [$P12PASS] output pem [$PEM_CERT_NAME] key [$PEM_KEY_NAME]"
openssl pkcs12 -clcerts -nokeys -out $PEM_CERT_NAME -in $P12CERT -passin pass:$P12PASS
openssl pkcs12 -nocerts -out $PEM_KEY_ENC_NAME -in $P12CERT -passin pass:$P12PASS -passout pass:$TMP_PASS
openssl rsa -in $PEM_KEY_ENC_NAME -out $PEM_KEY_NAME -passin pass:$TMP_PASS

echo "clean tmplate file '$PEM_KEY_ENC_NAME'"
rm  $PEM_KEY_ENC_NAME

echo "$PEM_CERT_NAME generate $PEM_CERT_NAME and $PEM_KEY_NAME ok"
