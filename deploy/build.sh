#!/bin/bash
set -euxo pipefail
cd `dirname $0`/

# use your own repo  
repository=michaelcourcy

# chart version
number=$(cat version.txt)
number=$((number+1))
echo $number > version.txt
version="0.0.$number"

# enter cmd directory to build images 
cd ../cmd/audit
GOOS=linux GOARCH=amd64 go build
docker build --platform=linux/amd64 -t $repository/audit-tool:$version-amd64 .
docker push $repository/audit-tool:$version-amd64
rm audit

cd `dirname $0`/
cat<<EOF > job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: audit-tool
  namespace: kasten-io
spec:
  template:
    spec:
      serviceAccount: k10-k10
      serviceAccountName: k10-k10
      containers:
      - name: audit-tool
        image: $repository/audit-tool:$version-amd64
        command: ["/audit"]
        env: 
        - name: KASTEN_NAMESPACE 
          value: kasten-io
        - name: KASTEN_RELEASE
          value: k10
      restartPolicy: Never
  backoffLimit: 4
EOF
if ! kubectl delete -n kasten-io -f job.yaml; 
then 
    echo "the audit-tools job was not there"
fi
kubectl create -n kasten-io -f job.yaml
sleep 20
kubectl logs -n kasten-io -f job/audit-tool