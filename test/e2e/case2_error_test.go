// Copyright (c) 2020 Red Hat, Inc.

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/open-cluster-management/governance-policy-propagator/test/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test error handling", func() {
	It("should not override remediationAction if doesn't exist on parent policy", func() {
		By("Creating ../resources/case2_error_test/remediation-action-not-exists.yaml on managed cluster in ns:" + testNamespace)
		utils.Kubectl("apply", "-f", "../resources/case2_error_test/remediation-action-not-exists.yaml",
			"-n", testNamespace)
		Eventually(func() interface{} {
			trustedPlc := utils.GetWithTimeout(clientManagedDynamic, gvrTrustedContainerPolicy,
				"case2-remedation-action-not-exsits-trustedcontainerpolicy", testNamespace, true, defaultTimeoutSeconds)
			return trustedPlc.Object["spec"].(map[string]interface{})["remediationAction"]
		}, defaultTimeoutSeconds, 1).Should(Equal("inform"))
		By("Patching ../resources/case2_error_test/remediation-action-not-exists2.yaml on managed cluster in ns:" + testNamespace)
		utils.Kubectl("apply", "-f", "../resources/case2_error_test/remediation-action-not-exists2.yaml",
			"-n", testNamespace)
		By("Checking the case2-remedation-action-not-exsits-trustedcontainerpolicy CR")
		yamlTrustedPlc := utils.ParseYaml("../resources/case2_error_test/remedation-action-not-exsits-trustedcontainerpolicy.yaml")
		Eventually(func() interface{} {
			trustedPlc := utils.GetWithTimeout(clientManagedDynamic, gvrTrustedContainerPolicy,
				"case2-remedation-action-not-exsits-trustedcontainerpolicy", testNamespace, true, defaultTimeoutSeconds)
			return trustedPlc.Object["spec"]
		}, defaultTimeoutSeconds, 1).Should(utils.SemanticEqual(yamlTrustedPlc.Object["spec"]))
		By("Deleting ../resources/case2_error_test/remediation-action-not-exists.yaml to clean up")
		utils.Kubectl("delete", "-f", "../resources/case2_error_test/remediation-action-not-exists.yaml",
			"-n", testNamespace)
		utils.ListWithTimeout(clientManagedDynamic, gvrPolicy, metav1.ListOptions{}, 0, true, defaultTimeoutSeconds)
	})
	It("should not break if no spec", func() {
		By("Creating ../resources/case2_error_test/no-spec.yaml on managed cluster in ns:" + testNamespace)
		utils.Kubectl("apply", "-f", "../resources/case2_error_test/no-spec.yaml",
			"-n", testNamespace)
		utils.GetWithTimeout(clientManagedDynamic, gvrPolicy, "default.case2-no-spec", testNamespace, true, defaultTimeoutSeconds)
		By("Deleting ../resources/case2_error_test/no-spec.yam to clean up")
		utils.Kubectl("delete", "-f", "../resources/case2_error_test/no-spec.yaml",
			"-n", testNamespace)
		utils.ListWithTimeout(clientManagedDynamic, gvrPolicy, metav1.ListOptions{}, 0, true, defaultTimeoutSeconds)
	})
})
