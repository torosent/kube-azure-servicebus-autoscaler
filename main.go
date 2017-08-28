package main

import (
	"flag"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/torosent/kube-azure-servicebus-autoscaler/azureservicebus"
	"github.com/torosent/kube-azure-servicebus-autoscaler/scale"
)

var (
	pollInterval             time.Duration
	scaleDownCoolPeriod      time.Duration
	scaleUpCoolPeriod        time.Duration
	scaleUpMessages          int
	scaleDownMessages        int
	maxPods                  int
	minPods                  int
	kubernetesDeploymentName string
	kubernetesNamespace      string
	resourcegroup            string
	queuename                string
	namespace                string
)

func Run(p *scale.PodAutoScaler) {
	lastScaleUpTime := time.Now()
	lastScaleDownTime := time.Now()

	for {
		select {
		case <-time.After(pollInterval):
			{
				numMessages, err := azureservicebus.NumMessages(resourcegroup, queuename, namespace)
				if err != nil {
					log.Errorf("Failed to get Azure Service Bus Queue messages: %v", err)
					continue
				}

				if numMessages >= scaleUpMessages {
					if lastScaleUpTime.Add(scaleUpCoolPeriod).After(time.Now()) {
						log.Info("Waiting for cool down, skipping scale up ")
						continue
					}

					if err := p.ScaleUp(); err != nil {
						log.Errorf("Failed scaling up: %v", err)
						continue
					}

					lastScaleUpTime = time.Now()
				}

				if numMessages <= scaleDownMessages {
					if lastScaleDownTime.Add(scaleDownCoolPeriod).After(time.Now()) {
						log.Info("Waiting for cool down, skipping scale down")
						continue
					}

					if err := p.ScaleDown(); err != nil {
						log.Errorf("Failed scaling down: %v", err)
						continue
					}

					lastScaleDownTime = time.Now()
				}
			}
		}
	}

}

func main() {
	flag.DurationVar(&pollInterval, "poll-period", 5*time.Second, "The interval in seconds for checking if scaling is required")
	flag.DurationVar(&scaleDownCoolPeriod, "scale-down-cool-down", 30*time.Second, "The cool down period for scaling down")
	flag.DurationVar(&scaleUpCoolPeriod, "scale-up-cool-down", 10*time.Second, "The cool down period for scaling up")
	flag.IntVar(&scaleUpMessages, "scale-up-messages", 100, "Number of azure service bus queue messages queued up required for scaling up")
	flag.IntVar(&scaleDownMessages, "scale-down-messages", 10, "Number of messages required to scale down")
	flag.IntVar(&maxPods, "max-pods", 5, "Max pods that kube-azure-servicebus-autoscaler can scale")
	flag.IntVar(&minPods, "min-pods", 1, "Min pods that kube-azure-service-autoscaler can scale")
	flag.StringVar(&kubernetesDeploymentName, "kubernetes-deployment", "", "Kubernetes Deployment to scale. This field is required")
	flag.StringVar(&kubernetesNamespace, "kubernetes-namespace", "default", "The namespace your deployment is running in")

	flag.StringVar(&resourcegroup, "resourcegroup", "", "Azure Service Bus resource group")
	flag.StringVar(&queuename, "queuename", "", "Azure Service Bus queue name")
	flag.StringVar(&namespace, "namespace", "", "Azure Service Bus namespace")

	flag.Parse()

	p := scale.NewPodAutoScaler(kubernetesDeploymentName, kubernetesNamespace, maxPods, minPods)
	log.Info("Starting kube-azure-service-autoscaler")
	Run(p)
}
