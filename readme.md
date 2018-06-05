# Traffic Monkey

Web application designed to be deployed in various kubernetes clusters from where to start send traffic to an endpoint.

urlHandlers:
 - /probe          - > send a get request and the monkey will start the attack
 - /data           - > http server to report data directory (exposes report files)
 - /.well-known/*  - > live / ready / metrics (no metrics)



## Current deployment process

In order to make it work, you need 2 terminal windows and to run a lot of commands :)

Warning: if --rm flag (default) is set and all tty's to the kubernetes pod are interrupted the pod will be destroyed and with him all the report data


```
Terminal 1
step 1 squarectl kubectl-shell -c mcc -v ccms && kubectl config get-contexts && kubectl config use-context <context>
step 2 kubectl run sre-shell --rm -i --tty --image ubuntu:18.04 -- bash
step 3 apt update && apt install -y curl
step 7 ./traefikmonkey

Terminal 2
step 4 squarectl kubectl-shell -c mcc -v ccms && kubectl config get-contexts && kubectl config use-context <context> && kubectl get po
step 5 kubectl cp ./trafficmonkey-linux-amd64 sre-shell-300457757-xk08h:/trafficmonkey
step 6 kubectl exec -i -t sre-shell-300457757-xk08h bash
step 8 curl 'http://localhost:9090/probe?url=https://idam.metrosystems.net/.well-known/openid-configuration&requests=1000&workers=10'
step 8 curl 'http://localhost:9090/probe?url=http://proxy.identity-prod:80/.well-known/openid-configuration&requests=1000&workers=10'
step 8 curl 'http://localhost:9090/probe?url=http://proxy-k8s-001-live1-mcc-gb-lon1.metroscales.io:30021/.well-known/openid-configuration&requests=1000&workers=10'
step 9 kubectl cp sre-shell-300457757-xk08h:/data/ ./data

```

## analyse generated data.

it was designed to generate percentile statistics and not detect 503 ( application code for timeout error ) :)

### analyse-data.sh
```
find . -type f -name "*.json" | while read file; do \
printf " >>> file: %s\n" $file; \
jq '."url"' $file; \
jq '."timestamp"' $file ; \
jq '."stats"' $file; \
jq '.["data"]' $file | grep status | sort | uniq -c; done
```

### Improvements

- implement gopkg.in/alecthomas/kingpin.v2
- implement github.com/sirupsen/logrus

3. Probe handler improvement: args insecure, debug
4. app entrypoint improvement: args insecure, debug
5. logrus debug
6. metrics  
7. save report by posting it to slack

```
http://localhost:9090/probe?insecure=true&resolve=10.29.30.8:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=4

http://localhost:9090/probe?url=http://localhost:9090&requests=10&workers=4
```