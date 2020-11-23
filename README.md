# Ledgerdata Refiner
Ledgerdata Refiner, proposed in [Ledgerdata Refiner: A Powerful Ledger Data Query Platform for Hyperledger Fabric](https://ieeexplore.ieee.org/abstract/document/8939212/), is a ledger data query platform for Hyperledger Fabric.  Its key component is ledger data analysis middleware, 
which extracts and synchronizes ledger data, and then parses the relationship among them. With this middleware, 
besides query blocks and transactions, we also provide enriched data view for end users, including schema overview and customized fine-grained query on ledger states.

The main features of Ledgerdata Refiner are as follows:

+ More flexible block and transaction retrieval and detail display
+ history tracking for all states
+ multiple chaincodes support
+ schema analysis for states in json format
+ fine-grained search for values in json format of states
+ A Dashboard to display system overview
![Fine grained query](https://github.com/FujitsuLaboratories/Ledger-Data-Refiner/blob/main/DemoUI/Advanced%20Query%20for%20Ledger%20Data%E2%80%94Ledger%20Data%20Refiner(1).png)

## Minimum requirements

| Requirements       | Notes          |
| ------------------ | -------------- |
| Go                 | 1.13 or higher |
| Hyperledger Fabric | 1.4.x and 2.x          |
| PostgreSQL         | 9.5 or higher  |

## Install

~~~shell
go get github.com/FujitsuLaboratories/Ledger-Data-Refiner
~~~



## Quick start(codebase)

### How to build hyperledger fabric network

See https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html.

### Database configuration

Modify `config/config.ini` to update PostgreSQL database settings:

~~~ini
[database]
db_host = 127.0.0.1
db_port = 5432
db_user = refiner
db_password = 123456
db_name = ledgerdata_refiner
~~~

### Fabric configuration (base on 2.2.x)

Modify `config/connection-config.yaml` to define your fabric network connection profile:

~~~yaml
client:
  organization: org1
  logging:
    level: info
  # Root of the MSP directories with keys and certs.
  cryptoconfig:
    path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}

organizations:
  org1:
    users:
      Admin:
        cert:
          path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/cert.pem
      User1:
        cert:
          path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem


orderers:
  orderer.example.com:
    tlsCACerts:
      # Certificate location absolute path
      path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
..............
peers:
  peer0.org1.example.com:
    tlsCACerts:
      # Certificate location absolute path
      path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
..............
ertificateAuthorities:
  ca.org1.example.com:
    url: https://ca.org1.example.com:7054
    tlsCACerts:
      # Comma-Separated list of paths
      path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
      # Client key and cert for SSL handshake with Fabric CA
      client:
        key:
          path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/tls.example.com/users/User1@tls.example.com/tls/client.key
        cert:
          path: ${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}/peerOrganizations/tls.example.com/users/User1@tls.example.com/tls/client.crt
............
~~~

Modify `${FABRIC_PROJECT_PATH}/${CRYPTOCONFIG_PATH}` to your path of keystore. For more information, please read [config_e2e.yaml](https://github.com/hyperledger/fabric-sdk-go/blob/master/test/fixtures/config/config_e2e.yaml)

### Run database script

If your database is deployed locally, run the following script to create user, database and tables:

~~~shell
$ cd DOCKER/postgreSQL/db
$ ./createdb.sh
~~~

Before running the script, you have to modify `config.json`:

~~~json
{
  "postgreSQL": {
    "database": "ledgerdata_refiner",
    "username": "refiner",
    "password": "123456"
  }
}
~~~

if you want to run the database on the docker, run the following command:

~~~shell
$ make docker_refinerdb
~~~

NOTE: the ledgerdata refiner uses [gorm](https://github.com/go-gorm/gorm) to connect to the database, it will automatically create the data table when the application starts, so you can create only user and database without tables.

### Build Ledgerdata Refiner

Use Makefile to build the application:

~~~shell
$ make build
~~~

or

~~~shell
$ go build -o ledgerdata-refiner main.go
~~~

### Run Ledgerdata Refiner

To run the application locally, run:

~~~shell
$ ./ledgerdata-refiner
~~~

Then open http://localhost:30052 in the browser



## Quick start(Docker)

### Configuration

1. Modify `config/config.ini` to update postgreSQL database settings:

   ~~~ini
   [database]
   # if using docker, db_host must be the same as the hostname in the docker-compose file
   db_host = refinerdb.network.com
   db_port = 5432
   db_user = refiner
   db_password = 123456
   db_name = ledgerdata_refiner
   ~~~

2. Modify `DOCKER/postgreSQL/db/config.json` to update postgreSQL database settings of docker image：

   ~~~json
   {
     "postgreSQL": {
       "database": "ledgerdata_refiner",
       "username": "refiner",
       "password": "123456"
     }
   }
   ~~~
   
3. Modify `docker-compose.yaml` to update volumes:

   ~~~yaml
   ..........
   
   refinerdb.network.com:
       volumes:
         - /you path here:/var/lib/postgresql/data
         
   ..........
   
   refiner.network.com:
       volumes:
         - /the path of your keystore:/opt/organizations
   .............
   ~~~

4. Build docker images:

   ~~~shell
   $ make docker_all
   ~~~

5. Run ledgerdata refiner:

   ~~~shell
   $ docker-compose -f docker-compose.yaml up -d
   ~~~

Then open http://localhost:30052 in the browser
## API document

The ledgerdata refiner uses swagger to provide API document. Once starting the application, check http://localhost:30052/swagger/index.html.

## Contributing

Thank you for considering to help out with the source code! We welcome contributions from anyone on the internet, and are grateful for even the smallest of fixes!
There are essentially 3 different scenarios for contributors.

### Fujitsu Employee or Contractor as Contributor
If you are a Fujitsu employee or contractor, you essentially already signed a CLA as part of your on-boarding. Simply fill in the required details and indicate that you are associated with Fujitsu in the Corporate CLA ([CCLA.txt](CCLA.txt)) .

### Individual Contributor
If you are contributing as an individual, simply fill in the required details and agree to the Individual CLA ([ICLA.txt](ICLA.txt)), then submit it to us.

### Corporate Contributor
If you are contributing as an employee of another company or organization, you need to simply fill in the required details and indicate the company you are associated with and that a Corporate CLA ([CCLA.txt](CCLA.txt)) has been signed and submitted to us.

## License

This Ledgerdata Refiner is distributed under the [Apache License Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html) found in the [LICENSE](LICENSE) file.