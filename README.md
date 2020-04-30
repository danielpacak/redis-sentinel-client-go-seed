![build](https://github.com/danielpacak/redis-sentinel-client-go-seed/workflows/build/badge.svg)

# redis-sentinel-client-go-seed

A seed Go app which connects to Redis sentinel cluster

## Deploying Redis Sentinel cluster

To begin the process, you need to obtain the `values.yaml` file included in the Bitnami Redis chart:

```
$ curl -Lo values.yaml https://raw.githubusercontent.com/bitnami/charts/master/bitnami/redis/values.yaml
```

Open the `values.yaml` file and edit the "sentinel" section as shown below:

```yaml
## Use redis sentinel in the redis pod. This will disable the master and slave services and
## create one redis service with ports to the sentinel and the redis instances
sentinel:
  enabled: true
## Use password authentication
usePassword: false
```

Install the latest version of the chart using the `values.yaml` file as shown below:

```
$ helm repo add bitnami https://charts.bitnami.com/bitnami
$ kubectl create namespace redis
$ helm install redis bitnami/redis -n redis --values values.yaml
```

Run a Redis pod that you can use as a client:

```
$ kubectl run --namespace redis redis-client --rm --tty -i --restart='Never' \
  --image docker.io/bitnami/redis:5.0.9-debian-10-r0 -- bash
If you don't see a command prompt, try pressing enter.
I have no name!@redis-client:/$
```

Sentinel access:

```
I have no name!@redis-client:/$ redis-cli -h redis -p 26379
redis.redis:26379> SENTINEL masters
1)  1) "name"
    2) "mymaster"
    3) "ip"
    4) "172.17.0.12"
    5) "port"
    6) "6379"
...
redis.redis:26379> exit
I have no name!@redis-client:/$ redis-cli -h 172.17.0.12 -p 6379
172.17.0.12:6379> SET foo "bar"
OK
172.17.0.12:6379> exit
I have no name!@redis-client:/$ exit
exit
pod "redis-client" deleted
```

## Getting started

```
$ kubectl apply -f kube/seed.yaml
$ kubectl port-forward -n seed service/seed 8080:8080
$ curl -d '{"key": "foo", "value": "bar"}' -H 'Content-Type: application/json' http://localhost:8080/redis/key
$ curl -H 'Accept: application/json' http://localhost:8080/redis/keys
["foo"]
```

## References

### Redis commands

| Command              | Description                                                                        |
| -------------------- | ---------------------------------------------------------------------------------- |
| [INFO][command-info] | Returns information and statistics about the server                                |
| [ROLE][command-role] | Provides information on the role of a Redis instance in the context of replication |

### Links

1. [Deploy a Redis Sentinel Kubernetes cluster using Bitnami Helm charts](https://docs.bitnami.com/tutorials/deploy-redis-sentinel-production-cluster)
2. [Redis Sentinel Documentation](https://redis.io/topics/sentinel)

[command-info]: https://redis.io/commands/info
[command-role]: https://redis.io/commands/role
