# devops-deploy-practice

## k3d + Helm

1) Install k3d
2) create k3d cluster `k3d cluster create <name> <optional portforwarding>`
3) cd into devops-deploy-practice dir
4) run `helm install ddp chart/`
5) run `helm ls -A`
6) should show ddp release