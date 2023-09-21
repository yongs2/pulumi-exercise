# [Learn Pulumi](https://www.pulumi.com/learn/)

## [Pulumi Fundamentals](https://www.pulumi.com/learn/pulumi-fundamentals/)

### [Creating a Pulumi Project](https://www.pulumi.com/learn/pulumi-fundamentals/create-a-pulumi-project/)

```sh
cd /workspace
mkdir -p 02.my_first_app && cd 02.my_first_app
pulumi new yaml -y
pulumi plugin install resource docker
```

### [Creating Docker Images](https://www.pulumi.com/learn/pulumi-fundamentals/create-docker-images/)

- refer https://github.com/pulumi/tutorial-pulumi-fundamentals

```sh
pulumi up --yes
docker images | grep "tutorial-pulumi"
```

### [Configuring and Provisioning Containers](https://www.pulumi.com/learn/pulumi-fundamentals/configure-and-provision/)

- config

```sh
pulumi config set frontendPort 3001
pulumi config set backendPort 3000
pulumi config set mongoPort 27017
pulumi config set mongoHost mongodb://mongo:27017
pulumi config set database cart
pulumi config set nodeEnvironment development
pulumi config set protocol http://
pulumi config
pulumi config -j

pulumi up --yes
```

- request

```sh
curl --location --request POST 'http://localhost:3000/api/products' \
-H 'Content-Type: application/json' \
--data-raw '{
    "ratings": {
        "reviews": [],
        "total": 63,
        "avg": 5
    },
    "created": 1600979464567,
    "currency": {
        "id": "USD",
        "format": "$"
    },
    "sizes": [
        "M",
        "L"
    ],
    "category": "boba",
    "teaType": 2,
    "status": 1,
    "_id": "5f6d025008a1b6f0e5636bc7",
    "images": [
        {
            "src": "classic_boba.png"
        }
    ],
    "name": "My New Milk Tea",
    "price": 5,
    "description": "none",
    "productCode": "852542-107"
}'
```

- response

```json
{"status":"ok","data":{"product":{"ratings":{"reviews":[],"total":63,"avg":5},"created":1600979464567,"currency":{"id":"USD","format":"$"},"sizes":["M","L"],"category":"boba","teaType":2,"status":1,"_id":"5f6d025008a1b6f0e5636bc7","images":[{"_id":"65013dc7685dff001b5ae59a","src":"classic_boba.png"}],"name":"My New Milk Tea","price":5,"description":"none","productCode":"852542-107","__v":0}}}
```

- [web page](http://localhost:3001)

### Cleaning up

```sh
pulumi destroy --yes
```

- To delete the stack itself

```sh
pulumi stack rm dev --yes
```

## [Building with Pulumi](https://www.pulumi.com/learn/building-with-pulumi/)

### [Understanding Stacks](https://www.pulumi.com/learn/building-with-pulumi/understanding-stacks/)

```sh
pulumi stack init dev
pulumi stack init staging
pulumi stack ls
pulumi stack select dev
pulumi stack ls
```

### [Understanding Stack Outputs](https://www.pulumi.com/learn/building-with-pulumi/stack-outputs/)

```sh
pulumi up -y
pulumi stack output url
curl $(pulumi stack output url)
```

- configurations in the dev stack

```sh
pulumi stack select dev
pulumi config
```

- configurations in the staging stack

```sh
pulumi stack select staging

pulumi config set frontendPort 3002
pulumi config set backendPort 3000
pulumi config set mongoPort 27017
pulumi config set mongoHost mongodb://mongo:27017
pulumi config set database cart
pulumi config set nodeEnvironment development
pulumi config set protocol http://

cat Pulumi.staging.yaml

pulumi up -y
pulumi stack output url
```

### Cleaning up

- staging

```sh
pulumi stack select staging
pulumi destroy --yes
```

- dev

```sh
pulumi stack select dev
pulumi destroy --yes
```
