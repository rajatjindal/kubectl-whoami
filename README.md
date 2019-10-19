# kubectl-whoami

`kubectl-whoami` is a [kubectl plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that show the subject that's currently authenticated as.

This plugin has been tested to work with following auth types:

- Basic Auth
- Cert Admin Auth
- RBAC Token in Kubeconfig file
- Token provided from command line using `--token` flag.

## TODO

- we want to test and add support for other auth-providers
- add unit tests
