#!/usr/bin/env bash

echo -ne "\033[32m*** Deploying ingress-nginx\033[0m\n"

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.27.1/deploy/static/mandatory.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.28.0/deploy/static/provider/cloud-generic.yaml

echo -ne "\033[32m*** Waiting for ingress-nginx, CTRL-C to exit\033[0m\n"

kubectl get pods --all-namespaces -l app.kubernetes.io/name=ingress-nginx #--watch
