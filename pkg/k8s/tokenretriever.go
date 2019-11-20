package k8s

import (
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
)

//WhatCanI runs self subject review request
func WhatCanI(kubeclient kubernetes.Interface) error {
	sar := &authorizationv1.SelfSubjectRulesReview{
		Spec: authorizationv1.SelfSubjectRulesReviewSpec{
			Namespace: "default",
		},
	}

	_, err := kubeclient.AuthorizationV1().SelfSubjectRulesReviews().Create(sar)
	return err
}
