# redis-ha-seed

## Getting started

```
$ curl -Lo values.yaml https://raw.githubusercontent.com/bitnami/charts/master/bitnami/redis/values.yaml
```

```
$ helm repo add bitnami https://charts.bitnami.com/bitnami
$ kubectl create namespace redis
$ helm install redis bitnami/redis -n redis --values values.yaml
```

```
$ kubectl run --namespace redis redis-client --rm --tty -i --restart='Never' \
  --image docker.io/bitnami/redis:5.0.9-debian-10-r0 -- bash
If you don't see a command prompt, try pressing enter.
I have no name!@redis-client:/$ redis-cli -h redis.redis -p 26379
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

```
$ kubectl apply -f kube/seed.yaml
$ kubectl port-forward -n seed service/seed 8080:8080
$ curl localhost:8080
["foo"]
```

## References

1. https://docs.bitnami.com/tutorials/deploy-redis-sentinel-production-cluster/
