# ldh-dns

create make file --> done
create docker file
cut release --> done
check for coredns file update --> done
code reorg --> done
create env for domain --> done
add lint

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

debug tools
kubectl run -it --image=praqma/network-multitool nettools --restart=Never --namespace=default
