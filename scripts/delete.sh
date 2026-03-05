kubectl delete -f ../manifests/automation/config-map.yaml
kubectl delete -f ../manifests/automation/secrets.yaml
kubectl delete -f ../manifests/automation/postgres-deployment.yaml
kubectl delete -f ../manifests/automation/deployment.yaml
kubectl delete -f ../manifests/automation/service.yaml
kubectl delete -f ../manifests/automation/name-space.yaml