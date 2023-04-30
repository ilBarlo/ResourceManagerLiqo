# FlavourGenerator
FlavourGenerator is a tool used to collect the resources of all the node connected to the SuperNode through Liqo and to collect the resources available for a specific Flavour (which in this case is defined as the **Architecture of the node** plus the **Operating System**) through a RESTful API. The project is based on the official [Liqo Dashboard project](https://github.com/LucaRocco/liqo-public-peering-dashboard/tree/8c10686e161b5853313ee9e3837dd6d25451dc78). 

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Prerequisites
The following packages are required to run FlavourGenerator:

- Go 1.16 or later
- Liqo (version 1.8 or later)

## Usage
To run FlavourGenerator The server locally you need to execute the run command in a console where you previously exported a KUBECONFIG variable.

```bash
export SERVER_ADDR="x.x.x.x:y"
go run cmd/main.go
```

To see the available resources for a specific Flavour, send a GET request to the following URL:

```bash
http://localhost:8080/api/flavour/available_resources
```

## Built With
- Go
- Liqo
