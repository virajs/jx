package cmd

import (
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

const (
	fromNs                   = "from-namespace-powdered-water"
	toNs                     = "to-namespace-journal-entry"
	serviceNameInFromNs      = "service-a-is-for-angry"
	serviceNameDummyInFromNs = "service-p-is-polluted"
	serviceNameInToNs        = "service-b-is-for-berserk"
)

func TestServiceLinking(t *testing.T) {
	o := StepLinkServicesOptions{
		FromNamespace: fromNs,
		Includes:      []string{serviceNameInFromNs},
		Excludes:      []string{serviceNameDummyInFromNs},
		ToNamespace:   toNs}
	fromNspc := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fromNs,
		},
	}
	svcInFromNs := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceNameInFromNs,
			Namespace: fromNs,
		},
	}
	svcDummyInFromNs := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceNameDummyInFromNs,
			Namespace: fromNs,
		},
	}
	toNspc := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: toNs,
		},
	}
	svcInToNs := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceNameInToNs,
			Namespace: toNs,
		},
	}

	ConfigureTestOptionsWithResources(&o.CommonOptions,
		[]runtime.Object{fromNspc, toNspc, svcInFromNs, svcInToNs, svcDummyInFromNs},
		nil,
		gits.NewGitCLI(),
		helm.NewHelmCLI("helm", helm.V2, ""))

	err := o.Run()
	serviceList, _ := o.kubeClient.CoreV1().Services(toNs).List(metav1.ListOptions{})
	serviceNames := []string{""}
	for _, service := range serviceList.Items {
		serviceNames = append(serviceNames, service.Name)
	}
	assert.Contains(t, serviceNames, serviceNameInFromNs) //Check if service that was in include list got added
	assert.EqualValues(t, len(serviceNames), 3)           //Check if service that was in exclude list didn't get added
	assert.NoError(t, err)
}
