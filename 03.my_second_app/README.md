# [Learn Pulumi](https://www.pulumi.com/learn/)

## [Building with Pulumi](https://www.pulumi.com/learn/building-with-pulumi/)

### [Understanding Stack References](https://www.pulumi.com/learn/building-with-pulumi/stack-references/)

```sh
cd /workspace
mkdir -p 03.my_second_app && cd 03.my_second_app
pulumi new yaml -y
pulumi stack init staging
pulumi stack ls
```

- set org configuration varaible

```sh
pulumi config set org organization
pulumi config

pulumi up -y
```

- get my_first_app's output in staing

```sh
pulumi stack output
```

- clean up

```sh
pulumi stack select staging
pulumi destroy --yes
```

### [Working with Secrets](https://www.pulumi.com/learn/building-with-pulumi/secrets/)

```sh
cd /workspace
cd 02.my_first_app

pulumi stack select dev

pulumi config
pulumi config set mongoUsername admin
pulumi config set --secret mongoPassword S3cr37
pulumi config
cat Pulumi.dev.yaml
pulumi config set mongoHost mongo

pulumi stack output mongoPassword --show-secrets
```

- clean up

```sh
pulumi stack select dev
pulumi destroy --yes
```
