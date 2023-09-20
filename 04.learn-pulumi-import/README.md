# [Learn Pulumi](https://www.pulumi.com/learn/)

## [Building with Pulumi](https://www.pulumi.com/learn/building-with-pulumi/)

### [Migration and Imports](https://www.pulumi.com/learn/importing/)

```sh
cd /workspace
mkdir -p 04.learn-pulumi-import && cd 04.learn-pulumi-import
pulumi new yaml -y

CONTAINER_ID=$(docker ps -a --no-trunc | grep frontend | awk '{printf $1}')
echo $CONTAINER_ID

pulumi import docker:index/container:Container frontend-dev $CONTAINER_ID -y
```
