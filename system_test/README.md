# kubelogin/system_test

This is an automated test for verifying behavior of the plugin with a real Kubernetes cluster and OIDC provider.

This code is essentially a copy of https://github.com/int128/kubelogin/blob/master/system_test/README.md + a test to verify subject name using kubectl-whoami plugin