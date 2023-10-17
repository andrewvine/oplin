# Oplin

Oplin is a simple service for collecting and viewing Metadata. It is similar to (but not as feature rich) as [Marquez](https://github.com/MarquezProject/marquez]). Like Marquez it implements the [OpenLineage](https://openlineage.io/) specification.


## Installation

1. Create a [postgreql](https://www.postgresql.org) database and user. One could run the following for example:

```
CREATE DATABASE oplin;
CREATE USER oplin WITH PASSWORD '{password}';
ALTER DATABASE oplin OWNER TO oplin; 
```

2. Setup the environment variables (replacing the defaults with your own as needed):

```
export OPLIN_DB_HOST = localhost;
export OPLIN_DB_USER = oplin; 
export OPLIN_DB_PASSWORD = {password};
export OPLIN_DB_NAME = oplin;
export OPLIN_DB_PORT = 5432;
```