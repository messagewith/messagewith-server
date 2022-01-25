# Messagewith Server
Messagewith server written in [Go](https://go.dev/)

## Running development server

### Requirements
- [Go toolchain](https://go.dev/doc/install)
- [MongoDB Server](https://www.mongodb.com/try/download/community)


Inside root folder, create .env file containing:
```
    MESSAGEWITH_JWT_SECRET=<random_long_text_for_jwt_sign>
    MESSAGEWITH_DOMAIN=<server_domain>
    MESSAGEWITH_MOCKUP_IP_ADDRESS=<mockup_ip_address_in_dev_env>
    MESSAGEWITH_DATABASE_URI=<mongodb_connection_uri>
```
Finally you can run ```go run .```


## License
This packages is distributed under the [GPL-3.0 License](https://github.com/messagewith/messagewith/blob/main/LICENSE).
