# Oplin

Oplin is a simple service for collecting and viewing Metadata. It is similar to (but not as feature rich) as [Marquez](https://github.com/MarquezProject/marquez]). Like Marquez it implements the [OpenLineage](https://openlineage.io/) specification.


## Installation


1. Download and unpack the latest [release](https://github.com/andrewvine/oplin/releases):

```
wget https://github.com/andrewvine/oplin/releases/download/v0.1.0/oplin_Linux_x86_64.tar.gz
tar -xvzf oplin_Linux_x86_64.tar.gz
```

2. Create a [postgreql](https://www.postgresql.org) database and user. From `psql` one could run:

```
CREATE DATABASE oplin;
CREATE USER oplin WITH PASSWORD '{password}';
ALTER DATABASE oplin OWNER TO oplin; 
```

3. Run the application. Replacing the default options as needed:

```
./oplin -db_host localhost -db_name oplin -db_password {password} -db_port 5432 -db_user oplin -web_port=8080
```
