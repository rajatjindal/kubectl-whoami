package k8s

import (
	"fmt"
	"regexp"
	"strings"

	authenticationapi "k8s.io/api/authentication/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

//WhoAmI returns the current user/token subject
func WhoAmI(kubeclient kubernetes.Interface, token string, printGroup bool) (string, error) {
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

	if result.Status.Error != "" {
		return "", fmt.Errorf(result.Status.Error)
	}
	userGroupStr := result.Status.User.Username

	if printGroup {
		userGroupStr = fmt.Sprintf("User:\t%s\nGroups:\n\t%s", result.Status.User.Username,
			strings.Join(result.Status.User.Groups, "\n\t"))
		if len(result.Status.User.Extra["arn"]) > 0 {
			userGroupStr = userGroupStr + "\nARN:\n\t" + strings.Join(result.Status.User.Extra["arn"], "\n\t")
		}
	}

	return userGroupStr, nil
}

func getUsernameFromError(err error) string {
	re := regexp.MustCompile(`^.* User "(.*)" cannot .*$`)
	return re.ReplaceAllString(err.Error(), "$1")
}
