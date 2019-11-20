package authproviders

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

//GetToken gets the token from auth provider config
func GetToken(ap *clientcmdapi.AuthProviderConfig) (string, error) {
	switch ap.Name {
	case "gcp":
		return ap.Config["access-token"], nil
	case "azure":
		return ap.Config["access-token"], nil
	case "oidc":
		return ap.Config["id-token"], nil
	}

	return "", fmt.Errorf("unsupported auth provider %s", ap.Name)
}
