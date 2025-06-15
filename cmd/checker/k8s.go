package checker

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sPodChecker struct {
	Namespace     string
	LabelSelector []string
	Scheme        string
	Port          int32
	Endpoint      string
	TLSVerify     bool
	Image         string

	Common *HTTPCommon
}

type kubeconfig struct {
	config    *rest.Config
	incluster bool
}

func (k *K8sPodChecker) getKubeConfig() (*kubeconfig, error) {
	incluster := true
	config, err := rest.InClusterConfig()
	if err == nil {
		return &kubeconfig{config, incluster}, nil
	} else {
		k.Common.Logger.Debug().Err(err).Msg("did not find incluster kube config. trying ~/.kube/config now...")
	}

	var kubeconfigPath string
	if kubeconfig_env, ok := os.LookupEnv("KUBECONFIG"); ok {
		kubeconfigPath = kubeconfig_env
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfigPath = filepath.Join(home, ".kube/config")
	}
	_, err = os.Stat(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	incluster = false
	return &kubeconfig{config, incluster}, nil
}

func (k *K8sPodChecker) Check(ctx context.Context) error {
	kubeconfig, err := k.getKubeConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig.config)
	if err != nil {
		return err
	}

	k.Common.Logger.Info().Strs("labels", k.LabelSelector).Str("namespace", k.Namespace).Msg("listing pods")

	pods, err := clientset.CoreV1().Pods(k.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(k.LabelSelector, ","),
	})
	if err != nil {
		return err
	}

	k.Common.Logger.Info().Msgf("found %d pods with the selected filters", len(pods.Items))

	var latestPod *corev1.Pod
	for _, pod := range pods.Items {
		if k.Image != "" && k.Image != pod.Spec.Containers[0].Image {
			k.Common.Logger.Info().Msgf("found pod %s/%s but did not match on image", pod.Namespace, pod.Name)
			continue
		}

		if latestPod == nil || latestPod.CreationTimestamp.After(pod.CreationTimestamp.Time) {
			latestPod = &pod
		}
	}

	if latestPod == nil {
		return fmt.Errorf("no matching pod found in %s namespace with labels: %s", k.Namespace, k.LabelSelector)
	}

	var podIP string
	for {
		podIP = latestPod.Status.PodIP
		if podIP != "" {
			k.Common.Logger.Info().Msgf("pod %s has an IP assigned: %s", latestPod.Name, latestPod.Status.PodIP)
			break
		}

		jitterSeconds := rand.Intn(6) + 5
		err := waitWithJitter(ctx, jitterSeconds)
		if err != nil {
			return err
		}
	}

	port := k.Port
	if port == 0 {
		if len(latestPod.Spec.Containers) > 0 && len(latestPod.Spec.Containers[0].Ports) > 0 {
			port = latestPod.Spec.Containers[0].Ports[0].ContainerPort
		}
	}

	var targetURL string
	if port == 0 {
		targetURL = fmt.Sprintf("%s://%s%s", k.Scheme, podIP, k.Endpoint)
	} else {
		targetURL = fmt.Sprintf("%s://%s:%d%s", k.Scheme, podIP, port, k.Endpoint)
	}

	k.Common.Logger.Info().Str("pod", latestPod.Name).Str("ip", podIP).Str("target", targetURL).Msg("constructed the target url")

	if !kubeconfig.incluster {
		k.Common.Logger.Info().Msg("skipping the healthcheck as we are not incluster and pod IP is inaccessible")
		return nil
	}

	staticCheck := &StaticHTTPChecker{
		Upstream: targetURL,
		Common:   k.Common,
	}

	return staticCheck.Check(ctx)
}
