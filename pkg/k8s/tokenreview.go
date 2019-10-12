package k8s

import (
	"regexp"

	authenticationapi "k8s.io/api/authentication/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

//WhoAmI returns the current user/token subject
func WhoAmI(kubeclient kubernetes.Interface, token string) (string, error) {
	result, err := kubeclient.AuthenticationV1().TokenReviews().Create(&authenticationapi.TokenReview{
		Spec: authenticationapi.TokenReviewSpec{
			Token: token,
		},
	})

	if err != nil {
		if k8serrors.IsForbidden(err) {
			return getUsernameFromError(err), nil
		}

		return "", err
	}

	return result.Status.User.Username, nil
}

func getUsernameFromError(err error) string {
	re := regexp.MustCompile(`^.* User "(.*)" cannot .*$`)
	return re.ReplaceAllString(err.Error(), "$1")
}
