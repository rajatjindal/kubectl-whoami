# kubectl-whoami

`kubectl-whoami` is a [kubectl plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that show the subject that's currently authenticated as.

This plugin has been tested to work with following auth types:

- Basic Auth
- Cert Admin Auth
- RBAC Token in Kubeconfig file
- Token provided from command line using `--token` flag.
- oidc provider

# Usage

<details><summary>start the minikube cluster (skip if you are using an existing cluster) </summary>
<p>

```bash
➜  kubectl-whoami git:(master) minikube start
😄  minikube v1.1.1 on darwin (amd64)
💡  Tip: Use 'minikube start -p <name>' to create a new cluster, or 'minikube delete' to delete this one.
🏃  Re-using the currently running virtualbox VM for "minikube" ...
⌛  Waiting for SSH access ...
🐳  Configuring environment for Kubernetes v1.14.3 on Docker 18.09.6
🔄  Relaunching Kubernetes v1.14.3 using kubeadm ... 
⌛  Verifying: apiserver proxy etcd scheduler controller dns
🏄  Done! kubectl is now configured to use "minikube"

## Observe that it has two contexts. One using basic-auth (default) and other using cert-auth (minikube)
➜  kubectl-whoami git:(master) cat ~/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1RENDQWN5Z0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFqTVJBd0RnWURWUVFLRXdkck0zTXQKYjNKbk1ROHdEUVlEVlFRREV3WnJNM010WTJFd0hoY05NVGt3TmpBeE1URXhNakF3V2hjTk1qa3dOVEk1TVRFeApNakF3V2pBak1SQXdEZ1lEVlFRS0V3ZHJNM010YjNKbk1ROHdEUVlEVlFRREV3WnJNM010WTJFd2dnRWlNQTBHCkNTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFER2RkMWtNQkh1bVp0WlFPUkpVdnczVzNKazFKTmEKOU4wUjRQcXYxb3FUWDVNaytsSVN6cTJscDdqM3dCQTVFUTNHcmJEVC9aTENpcXlDdjlBcGZtcmRqOWVBUGdFegowQjU1d0NISTBEZU9pUHJPUytKVkw2UmJUTXJNa0dpMHBwWVBnYWphZnhQZU1YNzdkZ3k2VzNSTS95R2dWNHF2CktqSURnczBtTXQ2RmUxYS9IRHM2WVFQSDBsM3FodjRnMHBtZnFyZE5OcE5vTkltYWluVUdXVGU3Q1ppSmR5MXQKeEo4R0VucCt6K3NxdWs3S2xhNnBLQ1ZBNXFGVEc4L3hXdjZadnFpMWVqa3RmM01PdzZDN3piNG9Sb1drS1JSVwo5TXpUWU1yR05GNXluRlNKZ2NRSThHZEFVQUF2R3R2VkxQVjIyZHdCNFJWcFJWRTJUN1Z2SVErOUFnTUJBQUdqCkl6QWhNQTRHQTFVZER3RUIvd1FFQXdJQ3BEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01BMEdDU3FHU0liM0RRRUIKQ3dVQUE0SUJBUUI0MzZLMmVsaW5pb0pLZmJHZmRuamFXdXllei9VVmNpbFRHdmpaTlp1d1hZYTJZMEhsK1hpWAplNGZnbDQzMEw5eXF5ZHBUQStBVGNIdDBmampuSkswKyt4QU1OM1QxbnFqeVU1c2tyY3l4YThxaFBxZENWSC92ClY2b3l0K0hUVVlUc2RHQW1tY29ISTE4NDBxMDk0NFRla0pHUXVOV0tUVEVxZ1BLRWorNklGellzaE85ZDdTQWEKZUFCWHBvSzE4djQvK2N0ek5RcXFqQi96WDNxd1pUdHp5bUVsVXk2WEJwWmpIbm9uei8rbWF5Y2ViM1B3OTQwQgp5djlrd284NENBL1h3NmVKSll5TmdGNTI0aGwwZnRsMEZ4Zk9IeldkdE9BcUM5YUFGOEwzWFhZZCtCZUU1aXVLClJCSHVjVnNlVFF3MFFJalQza2FoSUo2eHRiRThrV29VCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://192.168.1.2:6443
  name: default
- cluster:
    certificate-authority: /Users/rajatjindal/.minikube/ca.crt
    server: https://192.168.99.104:8443
  name: minikube
contexts:
- context:
    cluster: default
    user: default
  name: default
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: default
  user:
    password: 294f6f0dcfbf6a67ddb3737ce095ace7
    username: admin
- name: minikube
  user:
    client-certificate: /Users/rajatjindal/.minikube/client.crt
    client-key: /Users/rajatjindal/.minikube/client.key
```

- Run `kubectl-whoami` without any context/user override
```bash
  ➜  kubectl-whoami git:(master) ./kubectl-whoami 
kubecfg:certauth:admin
```
</p>
</details>


#### Run `kubectl-whoami` with `default` context
```bash
➜  kubectl-whoami git:(master) ./kubectl-whoami --context default 
kubecfg:basicauth:admin                                                                                       
```
<details><summary>Get token for a service account from the cluster and use that to authenticate</summary>
<p>

```bash
➜  kubectl-whoami git:(master) kubectl get secret -n kube-system
NAME                                             TYPE                                  DATA   AGE
attachdetach-controller-token-2rdxm              kubernetes.io/service-account-token   3      34d
bootstrap-signer-token-l79rf                     kubernetes.io/service-account-token   3      34d
certificate-controller-token-kchx9               kubernetes.io/service-account-token   3      34d
clusterrole-aggregation-controller-token-b68nk   kubernetes.io/service-account-token   3      34d
coredns-token-wndvv                              kubernetes.io/service-account-token   3      34d
cronjob-controller-token-hxjq9                   kubernetes.io/service-account-token   3      34d
daemon-set-controller-token-6p9br                kubernetes.io/service-account-token   3      34d
default-token-ls5lw                              kubernetes.io/service-account-token   3      34d
deployment-controller-token-9qj9k                kubernetes.io/service-account-token   3      34d
disruption-controller-token-7zsnk                kubernetes.io/service-account-token   3      34d
endpoint-controller-token-x2cd8                  kubernetes.io/service-account-token   3      34d
expand-controller-token-wpqh7                    kubernetes.io/service-account-token   3      34d
generic-garbage-collector-token-6n4p9            kubernetes.io/service-account-token   3      34d
horizontal-pod-autoscaler-token-qrmws            kubernetes.io/service-account-token   3      34d
job-controller-token-p9d7b                       kubernetes.io/service-account-token   3      34d
kube-proxy-token-9wlqp                           kubernetes.io/service-account-token   3      34d
namespace-controller-token-nfxnl                 kubernetes.io/service-account-token   3      34d
node-controller-token-44blg                      kubernetes.io/service-account-token   3      34d
persistent-volume-binder-token-kftqn             kubernetes.io/service-account-token   3      34d
pod-garbage-collector-token-d58dn                kubernetes.io/service-account-token   3      34d
pv-protection-controller-token-mqq2t             kubernetes.io/service-account-token   3      34d
pvc-protection-controller-token-b4c45            kubernetes.io/service-account-token   3      34d
replicaset-controller-token-4g52b                kubernetes.io/service-account-token   3      34d
replication-controller-token-59q77               kubernetes.io/service-account-token   3      34d
resourcequota-controller-token-sdjcs             kubernetes.io/service-account-token   3      34d
service-account-controller-token-pn7bk           kubernetes.io/service-account-token   3      34d
service-controller-token-d2gh7                   kubernetes.io/service-account-token   3      34d
statefulset-controller-token-hx4cb               kubernetes.io/service-account-token   3      34d
storage-provisioner-token-lml77                  kubernetes.io/service-account-token   3      34d
token-cleaner-token-fr7np                        kubernetes.io/service-account-token   3      34d
ttl-controller-token-7ntll                       kubernetes.io/service-account-token   3      34d
xyz                                              Opaque                                4      27d
➜  kubectl-whoami git:(master) kubectl get secret kube-proxy-token-9wlqp -o yaml -n kube-system
apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNU1ERXhPVEF6TkRBeU9Wb1hEVEk1TURFeE56QXpOREF5T1Zvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS3JGCjlORk85bzJnUG9wd2Uza3pxTmsyZUhvc1A3cU5TZHZqYVQ2bjJkTFVhVFZMSUUrbTlaM0JmUWZGbU5jSTE5K04KUWFkd1hha1I0ejAyVVlBd0ZWTjRIZHFLZll1MzNIWkNhdmZpVXE0ZE14ZWxVUWh4ZTBrM0Rvc216eFdSNExySwptbVNRTHhHb0ZyZHlPQzRhdG5xWnprVGJXWFdlU2VyL1F4elA2TUE1cWVEeGNXNlNEK2YyYi9NUCtJYUxLNHJaCk5tS3RHOVd1ZWhqRVhTWWhxVWFra3F1OVdaRUNGRlFHZUIxUkFLYUdtK0tlK2VBd004REpMSEVZRVU2aUkzdGgKY05haVVrNDlPK253T1kwT0VSa0EyVmQvN1lYUXZRaUZpSmRmdU1LTjE3eHpCYW1jMGtCU1dpeWxGR25lZG8zVwowODNTelhsN2U5NGhVR1NkU3FjQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFCK3ZaYzFRb25UWnIzOHVac2dDZGR2YlZ0TTZlN0FtMEQxK3lNUkZuTmM3cldqYmhZZgpuc0p0RE55ZnpzdFY1bzlUbTZTeUdDT3Z3aUIvbmQyL0MvRnVJQUd3TVdiQ0hxcG1FYnZRbEJwbDZWbGpTQmsrCm8rNCtMY3JHSVRNeDR0TDczU2k4UzhrejhWczUvRkxBVENxaThKRUhjblowNHhmcm40YUJlWUJKZTNtckJwMDcKbGw1SEdKS2NFcFYvTFRTdG03clAzTFVxMWRVMWRMYiswdFo3LzM1cGJ0RWYwNDM1QzBkalJVRmQ5NlFVQXlBMQpmeTcySHplLzhXOXR2REkrNXdORHZYOFRmRndEKzRIa1MrL1FLSDVoY0lucXNPdW5DR1lUbTdHR1N6RmY4Wnc5Cmp0YWo5cGQydnBMbFRXa2tQaDIyR3l2ZlJURGdPbFdpTy9HdgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  namespace: a3ViZS1zeXN0ZW0=
  token: ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSnJkV0psTFhONWMzUmxiU0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVmpjbVYwTG01aGJXVWlPaUpyZFdKbExYQnliM2g1TFhSdmEyVnVMVGwzYkhGd0lpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltdDFZbVV0Y0hKdmVIa2lMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzUxYVdRaU9pSmhZVGM1TldaaVl5MWtOekF5TFRFeFpUa3RZVE13TVMwd09EQXdNamRrWWpaalpHUWlMQ0p6ZFdJaU9pSnplWE4wWlcwNmMyVnlkbWxqWldGalkyOTFiblE2YTNWaVpTMXplWE4wWlcwNmEzVmlaUzF3Y205NGVTSjkuSXd4S2o0VUhfeldpTjlWSW9YTVhJOXlka19idHQxWUFLOFFzVkFET0Z4cmVoN01pa3BJcFNNN1NPUFdNZHN3VjZRd1FjbnZScjBGQWNmTnBVOE5VM0lrSGIybllkOFhDWm1iZDhaeEliRm1Velc1VVRjdGZpa0ZHMmxJN1V5MngwNmg3RXMzWWFyMEhuZHJaZi01UWg3cUQxQkhtUHZ5MkxxTS1vYVNHamtkLVRLd1l2T3VTb2I3OUxZczdMdGRCeFA5UGRlTXZZUDd4RklBU3VmZzczSDVERlRxV3ZYS0tDemNIZHNoUTE0QXhtR3Z2dDlreWdSeFVKSFBmNGxXbU9rTWVldjNIUHRiMDNSN19hRkx3YldnTzJnTVFQTFlyaFpEc3laYVdyUVljdnhkaFhWSXUwN0lCOGJBN1F0OExqYnllX3V6RDJJRVdQaktZYjFnQ2dR
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: kube-proxy
    kubernetes.io/service-account.uid: aa795fbc-d702-11e9-a301-080027db6cdd
  creationTimestamp: "2019-09-14T15:16:46Z"
  name: kube-proxy-token-9wlqp
  namespace: kube-system
  resourceVersion: "213"
  selfLink: /api/v1/namespaces/kube-system/secrets/kube-proxy-token-9wlqp
  uid: aa7a8f87-d702-11e9-a301-080027db6cdd
type: kubernetes.io/service-account-token
➜  kubectl-whoami git:(master) echo ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSnJkV0psTFhONWMzUmxiU0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVmpjbVYwTG01aGJXVWlPaUpyZFdKbExYQnliM2g1TFhSdmEyVnVMVGwzYkhGd0lpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltdDFZbVV0Y0hKdmVIa2lMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzUxYVdRaU9pSmhZVGM1TldaaVl5MWtOekF5TFRFeFpUa3RZVE13TVMwd09EQXdNamRrWWpaalpHUWlMQ0p6ZFdJaU9pSnplWE4wWlcwNmMyVnlkbWxqWldGalkyOTFiblE2YTNWaVpTMXplWE4wWlcwNmEzVmlaUzF3Y205NGVTSjkuSXd4S2o0VUhfeldpTjlWSW9YTVhJOXlka19idHQxWUFLOFFzVkFET0Z4cmVoN01pa3BJcFNNN1NPUFdNZHN3VjZRd1FjbnZScjBGQWNmTnBVOE5VM0lrSGIybllkOFhDWm1iZDhaeEliRm1Velc1VVRjdGZpa0ZHMmxJN1V5MngwNmg3RXMzWWFyMEhuZHJaZi01UWg3cUQxQkhtUHZ5MkxxTS1vYVNHamtkLVRLd1l2T3VTb2I3OUxZczdMdGRCeFA5UGRlTXZZUDd4RklBU3VmZzczSDVERlRxV3ZYS0tDemNIZHNoUTE0QXhtR3Z2dDlreWdSeFVKSFBmNGxXbU9rTWVldjNIUHRiMDNSN19hRkx3YldnTzJnTVFQTFlyaFpEc3laYVdyUVljdnhkaFhWSXUwN0lCOGJBN1F0OExqYnllX3V6RDJJRVdQaktZYjFnQ2dR | base64 --decode

eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlLXByb3h5LXRva2VuLTl3bHFwIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Imt1YmUtcHJveHkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJhYTc5NWZiYy1kNzAyLTExZTktYTMwMS0wODAwMjdkYjZjZGQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06a3ViZS1wcm94eSJ9.IwxKj4UH_zWiN9VIoXMXI9ydk_btt1YAK8QsVADOFxreh7MikpIpSM7SOPWMdswV6QwQcnvRr0FAcfNpU8NU3IkHb2nYd8XCZmbd8ZxIbFmUzW5UTctfikFG2lI7Uy2x06h7Es3Yar0HndrZf-5Qh7qD1BHmPvy2LqM-oaSGjkd-TKwYvOuSob79LYs7LtdBxP9PdeMvYP7xFIASufg73H5DFTqWvXKKCzcHdshQ14AxmGvvt9kygRxUJHPf4lWmOkMeev3HPtb03R7_aFLwbWgO2gMQPLYrhZDsyZaWrQYcvxdhXVIu07IB8bA7Qt8Ljbye_uzD2IEWPjKYb1gCgQ%                  
```

</p>
</details>

#### use the token at command line to get its subject

```bash
➜  kubectl-whoami git:(master) ./kubectl-whoami --token eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlLXByb3h5LXRva2VuLTl3bHFwIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Imt1YmUtcHJveHkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJhYTc5NWZiYy1kNzAyLTExZTktYTMwMS0wODAwMjdkYjZjZGQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06a3ViZS1wcm94eSJ9.IwxKj4UH_zWiN9VIoXMXI9ydk_btt1YAK8QsVADOFxreh7MikpIpSM7SOPWMdswV6QwQcnvRr0FAcfNpU8NU3IkHb2nYd8XCZmbd8ZxIbFmUzW5UTctfikFG2lI7Uy2x06h7Es3Yar0HndrZf-5Qh7qD1BHmPvy2LqM-oaSGjkd-TKwYvOuSob79LYs7LtdBxP9PdeMvYP7xFIASufg73H5DFTqWvXKKCzcHdshQ14AxmGvvt9kygRxUJHPf4lWmOkMeev3HPtb03R7_aFLwbWgO2gMQPLYrhZDsyZaWrQYcvxdhXVIu07IB8bA7Qt8Ljbye_uzD2IEWPjKYb1gCgQ
system:serviceaccount:kube-system:kube-proxy
➜ 
```

## TODO

- add unit tests
