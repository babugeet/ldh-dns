# ldh-dns

This is a simple DNS server cum Response re-writer (for kubernetes DNS) written for study purpose. The main purpose of the server is to achieve split-horizon for DNS queries. Most of the DNS servers/operator (core DNS) already have options to achieve split horizon. For the servers which doesnt have this capability can use something like this to achieve the same. If the endpoint is of public and if we want the traffic to be redirected to an internal service in the k8s cluster, this code can be used.

## Pre-req and Conditions

- This server will be listening on localhost 53
- Server requires two env to be set.
    - `_DNS_SERVER` external DNS IP need to be provided. This IP will be used to forward the unresolved DNS queries
    - `_PRIVATE_DOMAIN`, This domain name will be the triggering point for the split horizon
- The url which need to be intercepted by this DNS server need to be in the format `test.<namespace>._PRIVATE_DOMAIN`, `eg: test.namespace.sample.dev.ldhappdomain.cloud`
- The above url format will be resolved to  `test.<namespace>.svc.cluster.local`; if the url is of the format   `test.sample.dev.ldhappdomain.cloud`, then it will be redirected to a service in the `linuxdatahub` namespace, `eg: test.linuxdatahub.sample.dev.ldhappdomain.cloud`
- For the scope of this repo, a simple pod with the server code is ran and the same is exposed internally to the cluster via a ClusterIP service
- DNS config file need to be updated to include the new DNS server. Below is the sample file in case of coreDNS, wherein 10.105.54.175:53 is the endpoint where DNS is hosted.

```
  Corefile: |
    sample.dev.ldhappdomain.cloud:53 {
         prometheus 127.0.0.1:9153
         forward .  10.105.54.175:53 {
             policy random
         }
         errors
         log . {
             class error
         }
         bufsize 1232
         cache 900 {
             denial 9984 30
         }
     }
    .:53 {
        errors
        health {
            lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
            pods insecure
            fallthrough in-addr.arpa ip6.arpa
            ttl 30
        }
        rewrite {
          name regex (.*)\.nginx\.svc\.cluster\.local name.default.svc.cluster.local
          answer name namei1.default.svc.cluster.local  main1-nginx-ingress-controller.nginx.svc.cluster.local
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

## Debug tools
A simple container which will have all debug tools
```
kubectl create deploy utils --image=arunvelsriram/utils --replicas=1 -- sleep infinity

kubectl exec -it deploy/utils -- bash
```