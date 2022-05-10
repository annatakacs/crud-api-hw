### APPLICATION DESCRIPTION

CRUD REST API service implemented in golang.
Following features have been implemented:
  - initializing and migrating sample data into postgresql database
  - execute GET, PUT, POST, DELETE operations on postgresql database
  - dockerize go environment
  - deploy postgresql database and dockerized crud api service with helm

### How to start the service?

1) Please install minikube on your machine: https://minikube.sigs.k8s.io/docs/start/

2) Open up a Terminal on your machine and navigate to the application directory.

3) Initialize the minikube cluster with the following command:
```
$ minikube start --driver=hyperkit
```
4) Navigate to the helm directory:
```
$ cd ./helm
```

5) Execute the following commands to start up the services:

```
$ helm install -f meals-postgres.yaml postgres ./postgres
```

Example output of the command:
```
$ helm install -f meals-postgres.yaml postgres ./postgres
NAME: postgres
LAST DEPLOYED: Mon May  9 21:56:32 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

```
$ helm install crud ./crud
```

Example output of the command:
```
$ helm install crud ./crud
NAME: crud
LAST DEPLOYED: Mon May  9 21:56:42 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

4) Open a new Terminal and create a tunnel for the minikube cluster:
```
$ minikube tunnel
```

Example output of the command:
```
$ minikube tunnel
Password:
Status: 
    machine: minikube
    pid: 2556
    route: 10.96.0.0/12 -> 192.168.64.2
    minikube: Running
    services: [crud-api, hello-minikube1]
    errors: 
        minikube: no errors
        router: no errors
        loadbalancer emulator: no errors
```

Please note that it is important not to close this terminal window as long as you are using the CRUD API service.

5) Return to your previous terminal and retrieve the EXTERNAL-IP address of crud-api service for your queries with the following command:
```
$ kubectl get svc
```

Example output of the command:
```
$ kubectl get svc
NAME              TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)          AGE
crud-api          LoadBalancer   10.111.79.235   10.111.79.235   8080:31350/TCP   2m29s
postgres          ClusterIP      10.111.85.99    <none>          5432/TCP         2m39s
```

### Example queries:

GET all meals: `http://${crud-api-EXTERNAL-IP}:8080/api/v1/get/meals`
GET meal by ID: `http://${crud-api-EXTERNAL-IP}:8080/api/v1/get/meals/{id}`
POST meal: `http://${crud-api-EXTERNAL-IP}:8080/api/v1/post/meals`

Example post request body:
```
{
    "name": "Vegan burger",
    "price": 9.45,
    "ingredients": "mushrooms, onions, bun, ketchup, mayo",
    "spicy": false,
    "vegan": true,
    "glutenFree": false,
    "description": "Delicious vegan burger",
    "kcal": 450
}
```

Expected response:
```
"Succesfully inserted data"
```

PUT meal: `http://${crud-api-EXTERNAL-IP}:8080/api/v1/put/meals/{id}`

Example put request body:
```
{
    "name": "Vegan burger",
    "price": 9.45,
    "ingredients": "mushrooms, onions, bun, ketchup, mayo",
    "spicy": false,
    "vegan": true,
    "glutenFree": false,
    "description": "Delicious vegan burger",
    "kcal": 450
}
```

DELETE meal by ID: `http://${crud-api-EXTERNAL-IP}:8080/api/v1/delete/meals/{id}`

Expected response:

```
"Succesfully deleted row"
```

### Cleanup steps

1) Delete minikube cluster:
```
$ minikube delete
```