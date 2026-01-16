# Smart Gas Payments
---
## Backend service (API) that helps doing payments

![Smart gas image](https://smartgasgasolineras.mx/images/logos/Logo_SmartGas.png)

Service made with golang in order to provide fast and a big amount of requests per second.

### OS Environment

#### App General variables
| VAR NAME      | Description                                                   | Default   |
| -----------   | -----------                                                   | --------  | 
| PORT          | Application port                                              |  8008     |
| HOST          | What host will be running                                     |  localhost:8008  |
| DEBUG         | Is application on debug?                                      |           |
| SECRET_KEY    | Secret key used for sign credentials                          |           |
| SECRET_KEY_REFRESH | Secret key for the refresh token                         |           |
| JWT_EXP_MINUTES |  Jwt Epiration in minutes                                   |           |
| JWT_REFRESH_EXP_DAYS |  JWT REFRESH EXPIRATION IN DAYS                        |           |
| TZ            | *(Important)* Timezone that will be used on the timezone      | UTC (if not setted in the OS)|
| TRUSTED_PROXIES            | Allowed Trusted Proxies, example: google.com youtube.com      | * |
| ALLOWED_HOSTS            | Allowed hosts, example: google.com youtube.com      | * |
| SOCIO_SMART_URL            | SocioSmartUrl      |  |
| STRIPE_SECRET_KEY            | Stripe secret key      |  |
| STRIPE_WEBHOOK_SECRET            | Stripe secret key for webhook      |  |
| ENABLE_GAS_PUMP            | Enable gas pump      |  False |
| SWIT_BASE_URL            | Swit Base Url      |   |
| PAYMENT_PROVIDER            | Payment Provider      |  stripe |
| CONECTIA_URL            | Conectia Url      |   |
| CONECTIA_URL_API            | Conectia Url Api      |   |
| CONECTIA_TOKEN            | Conectia Token      |   |
| SMTP_USER            | Smtp user      |   |
| SMTP_HOST            | Smtp host      |   |
| SMTP_PORT            | Smtp Port      |   |
| SMTP_PASSWORD            | Smtp Password      |   |
| FROM_EMAIL            | Email used for email notifications      |   |
| SENTRY_DSN            | Sentry Dsn      |   |
| ENVIRONMENT            | Environment | development  |




#### Databas variables
| VAR NAME      | Description                                                   | Default   |
| -----------   | -----------                                                   | --------  | 
| DB_HOST       | Database Host                                                 |           |
| DB_NAME       | Database Name                                                 |           |
| DB_USER       | Database User                                                 |           |
| DB_PASS       | Database Password                                             |           |
| DB_PORT       | Database PORT                                                 |  3306     |


### How to compile the docker image

>$ docker build -t smartgas-payments-backend .

### How to run the project

>$ docker run -e PORT=8008 -e TZ=UTC ... -e DB_PORT=3306 smartgas-payments-backend

*(Important)* Populate all the neccesary variables





--- 
### TODOS:

#### Auth:
* [x] Login
* [x] Refresh token

#### Users:
* [x] List Users
* [x] Me
* [x] Create Users (admin)

#### Payments:
* [x] Make a payment in order to refuel tank
* [x] Create socket that notifies when a gas pump gas finished the proccess of load fuel

#### Wallet:
* [x] Create wallet if customer does not have
* [x] List wallet cards
* [x] Set default card


#### CLI
* [x] Run Server
* [x] Migrate database
* [x] Create User

#### Core
* [x] Project Structure
* [x] Env configuration
* [x] Dependency Injection


#### Documentation
* [x] Initialize Swagger
* [x] Document Application
* [x] Expose documentation
