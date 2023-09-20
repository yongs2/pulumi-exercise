import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as util from "./util";

// Install Prometheus on the cluster.
const prometheus = new k8s.helm.v3.Chart("p8s", {
    repo: "prometheus",
    chart: "prometheus",
    version: "21.1.2",
    values: {
        "prometheus-node-exporter": {
            service: {
                port: 9101
            }
        },
        alertmanager: {
            persistence: {
              storageClass : "openebs-hostpath"
            }
        },
        server: { 
            persistentVolume: { 
                storageClass: "openebs-hostpath"
            }
        }
    }
});

const containerName = "example-app";
const appLabels = { app: containerName };
// Define Pod template that deploys an instrumented app. Annotation `prometheus.io/scrape` instructs
// Prometheus to scrape this pod for metrics.
const instrumentedPod = {
    metadata: { 
        annotations: { 
            "prometheus.io/scrape": "true" 
        }, 
        labels: appLabels
    },
    spec: {
        containers: [
            {
                name: containerName,
                // Prometheus-instrumented app that generates artificial load on itself.
                image: "fabxc/instrumented_app",
                ports: [{ name: "web", containerPort: 8080 }],
            },
        ],
    },
};

// kubectl -n default get Service p8s-prometheus-server -o yaml
const p8sService = prometheus.getResource("v1/Service", "default/p8s-prometheus-server");
// p8sService.apply((p8sService) => console.log("p8sService: %o", p8sService));

// kubectl -n default get Deployment p8s-prometheus-server -o yaml
const p8sDeployment = prometheus.getResource("apps/v1/Deployment", "default/p8s-prometheus-server");
// p8sDeployment.apply((p8sDeployment) => console.log("p8sDeployment: %o", p8sDeployment));

// IMPORTANT: This forwards the Prometheus service to localhost, so we can check it. If you are
// running in-cluster, you probably don't need this!
const localPort = 9090;
const forwarderHandle = util.forwardPrometheusService(p8sService, p8sDeployment, {
    localPort,
});

// Canary ring. Replicate instrumented Pod 3 times.
const canary = new k8s.apps.v1.Deployment(
    "canary-example-app", {
        spec: { 
            replicas: 1,
            selector: {
                matchLabels: appLabels
            },
            template: instrumentedPod 
        } 
    },
    { 
        dependsOn: p8sDeployment 
    },
);

// Staging ring. Replicate instrumented Pod 10 times.
const staging = new k8s.apps.v1.Deployment("staging-example-app", {
    metadata: {
        annotations: {
            // Check P90 latency is < 100,000 microseconds. Returns a `Promise<string>` with the P90
            // response time. It must resolve correctly before this deployment rolls out. In
            // general any `Promise<T>` could go here.
            "example.com/p90ResponseTime": util.checkHttpLatency(canary, containerName, {
                durationSeconds: 60,
                quantile: 0.9,
                thresholdMicroseconds: 100000,
                prometheusEndpoint: `localhost:${localPort}`,
                forwarderHandle,
            }),
        },
    },
    spec: { 
        replicas: 1,
        selector: {
            matchLabels: appLabels
        },
        template: instrumentedPod 
    },
});

export const p90ResponseTime = staging.metadata.annotations["example.com/p90ResponseTime"];
