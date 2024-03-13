# Audit tool for Kasten 

## Goal 

Create an audit of your Kasten backup on your cluster. 

## Features

- Check Kubernetes installation 
- Check Kasten installation 
- Analysing profiles 
  - Detect no profiles and no location profiles 
  - Detect profile not valid 
  - Detect no profile with immutability
- Give a RPO for each of your namespaces starting by namespaces having PVC 
- Give a RPO for each of your namespaces starting by namespaces not having PVC 

Coming soon :
- Is disaster recovery activated 
- Is RBAC applied to let non admin user manage their backups/restore
- Is Garbage collector enabled 
- Are database workloads using blueprints 
- Licencing checking 

## Deploy 

You must have Kasten installed on your cluster. We assume that 
Kasten is installed in the kasten-io namespace under the release k10.

If not you can change those values in the file `deploy/job.yaml` you should 
also change the `serviceAccount` and `serviceAccountName` accordingly. 


Deploy the audit job to your cluster 
```
kubectl create -f deploy/job.yaml 
```

Then you can check the logs 
```
kubectl logs -n kasten-io -f job/audit-tool
```

You should have an output similar to this one 
```
Kasten is installed in kasten-io under the release k10 

======================
Runtime for audit tool
======================
Audit tool is executing in pod

===================
Namespaces with PVC
===================
xxxxx-vm2 has 1 PVCs
  --> No backupactions in namespace xxxxx-vm2
xxxxxx-pacman has 1 PVCs
  --> No backupactions in namespace xxxxxx-pacman 
account-management has 1 PVCs
  BACKUPACTION                   STATE      START                          STOP                           
  scheduled-r6xr6                Complete   2024-03-06 16:00:17 +0000 UTC  2024-03-06 16:01:54 +0000 UTC  
  scheduled-qwp7g                Complete   2024-03-06 15:00:09 +0000 UTC  2024-03-06 15:01:47 +0000 UTC  
  scheduled-cc582                Complete   2024-03-06 14:00:20 +0000 UTC  2024-03-06 14:02:11 +0000 UTC  
  scheduled-r64c8                Complete   2024-03-06 13:00:11 +0000 UTC  2024-03-06 13:01:55 +0000 UTC  
  scheduled-x5zm9                Complete   2024-03-06 12:00:18 +0000 UTC  2024-03-06 12:01:51 +0000 UTC  
  scheduled-6h2m7                Complete   2024-03-06 11:00:10 +0000 UTC  2024-03-06 11:01:48 +0000 UTC  
  scheduled-d9tsr                Complete   2024-03-06 10:00:14 +0000 UTC  2024-03-06 10:01:57 +0000 UTC  
  scheduled-bmr4d                Complete   2024-03-06 09:00:21 +0000 UTC  2024-03-06 09:01:43 +0000 UTC  
  scheduled-ngsrm                Complete   2024-03-06 08:00:13 +0000 UTC  2024-03-06 08:01:32 +0000 UTC  
  scheduled-x6sb8                Complete   2024-03-06 07:00:19 +0000 UTC  2024-03-06 07:01:41 +0000 UTC  
  scheduled-8jfmb                Complete   2024-03-06 06:00:10 +0000 UTC  2024-03-06 06:01:43 +0000 UTC  
  scheduled-bb2t7                Complete   2024-03-06 05:00:10 +0000 UTC  2024-03-06 05:01:29 +0000 UTC  
  scheduled-vwbwt                Complete   2024-03-06 04:00:16 +0000 UTC  2024-03-06 04:01:31 +0000 UTC  
  scheduled-sf64t                Complete   2024-03-06 03:00:09 +0000 UTC  2024-03-06 03:01:48 +0000 UTC  
  scheduled-xh7l9                Complete   2024-03-06 02:00:13 +0000 UTC  2024-03-06 02:01:31 +0000 UTC  
  scheduled-l7dk5                Complete   2024-03-06 01:00:20 +0000 UTC  2024-03-06 01:01:56 +0000 UTC  
  scheduled-jlpzd                Complete   2024-03-06 00:00:17 +0000 UTC  2024-03-06 00:02:00 +0000 UTC  
  scheduled-xp9xp                Complete   2024-03-05 23:00:09 +0000 UTC  2024-03-05 23:01:39 +0000 UTC  
  scheduled-bnqw2                Complete   2024-03-05 22:00:10 +0000 UTC  2024-03-05 22:01:47 +0000 UTC  
  scheduled-bvqq9                Complete   2024-03-05 21:00:15 +0000 UTC  2024-03-05 21:01:39 +0000 UTC  
  scheduled-h59qr                Complete   2024-03-05 20:00:20 +0000 UTC  2024-03-05 20:01:58 +0000 UTC  
  scheduled-rwxhz                Complete   2024-03-05 19:00:10 +0000 UTC  2024-03-05 19:01:36 +0000 UTC  
  scheduled-h26bv                Complete   2024-03-05 18:00:17 +0000 UTC  2024-03-05 18:01:49 +0000 UTC  
  scheduled-2md5g                Complete   2024-03-05 17:00:09 +0000 UTC  2024-03-05 17:01:42 +0000 UTC  
  scheduled-k89l4                Complete   2024-03-05 00:00:11 +0000 UTC  2024-03-05 00:01:49 +0000 UTC  
  scheduled-krl6k                Complete   2024-03-04 00:00:18 +0000 UTC  2024-03-04 00:01:41 +0000 UTC  
  scheduled-kc6s6                Complete   2024-03-03 00:00:22 +0000 UTC  2024-03-03 00:01:55 +0000 UTC  
  scheduled-smzgz                Complete   2024-03-02 00:00:14 +0000 UTC  2024-03-02 00:01:48 +0000 UTC  
  scheduled-xc8z9                Complete   2024-03-01 00:00:17 +0000 UTC  2024-03-01 00:01:36 +0000 UTC  
  scheduled-5dbxg                Complete   2024-02-29 00:00:11 +0000 UTC  2024-02-29 00:01:46 +0000 UTC  
  scheduled-c6rd9                Complete   2024-02-25 00:00:10 +0000 UTC  2024-02-25 00:01:37 +0000 UTC  
  scheduled-kl78j                Complete   2024-02-18 00:00:20 +0000 UTC  2024-02-18 00:01:57 +0000 UTC  
  scheduled-5rzhl                Complete   2024-02-11 00:00:12 +0000 UTC  2024-02-11 00:01:37 +0000 UTC  
  scheduled-g9fzv                Complete   2024-02-09 12:59:26 +0000 UTC  2024-02-09 13:00:41 +0000 UTC  
  scheduled-j9rvl                Complete   2024-02-09 10:47:59 +0000 UTC  2024-02-09 10:49:17 +0000 UTC  
  --> The last RPO is 0 days and 2 hours
broken-pieces has 1 PVCs
  BACKUPACTION                   STATE      START                          STOP                           
  scheduled-l6rrr                Failed     2024-01-05 18:45:57 +0000 UTC  2024-01-05 18:47:30 +0000 UTC  
  --> WARNING !! It seems that no backupaction were successful
test-fedora has 1 PVCs
  BACKUPACTION                   STATE      START                          STOP                           
  scheduled-tzcbt                Complete   2024-02-04 04:18:50 +0000 UTC  2024-02-04 04:19:57 +0000 UTC  
  --> The last RPO is 31 days and 11 hours
```

## Contribute 

We welcome PRs 

## Build 

you must have golang and docker installed.

Change the value of the `repository` variable in `deploy/build.sh` by your docker repository. 
You can also change the semeantic value of the `version` variable. 

then to deploy your changes and create a new deployable image : 
```
./deploy/build.sh
```

It will recreate, push and redeploy a new job with the new image.

## Local development 

Assuming your context point to your cluster 

```
cd cmd/audit 
KASTEN_NAMESPACE=kasten-io \
KASTEN_RELEASE=k10 \
go run main.go
```


