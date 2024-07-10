jwt_key:
	ssh-keygen -t rsa -b 4096 -m PEM -f banking_jwtRS256.key

jwt_pub:
	openssl rsa -in banking_jwtRS256.key -pubout -outform PEM -out banking_jwtRS256.key.pub

jwt: jwt_key jwt_pub