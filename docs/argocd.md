# ArgoCD

When installing using ArgoCD you will need to tell ArgoCD to ignore some changes.

The helm chart makes use of `genSignedCert` to configure and secure the webhooks.
Because of this implementation, each `helm template` will cause a change and will as a result 
keep creating differences for ArgoCD to deploy.

Luckily ArgoCD has a buildin mechanism to handle these kinds of workflows. The example below can be used 
to ignore these changes.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
spec:
  ignoreDifferences:
    - group: admissionregistration.k8s.io
      kind: ValidatingWebhookConfiguration
      jsonPointers:
        - /webhooks/0/clientConfig/caBundle
        - /webhooks/1/clientConfig/caBundle
    - group: admissionregistration.k8s.io
      kind: MutatingWebhookConfiguration
      jsonPointers:
        - /webhooks/0/clientConfig/caBundle
        - /webhooks/1/clientConfig/caBundle
    - group: ""
      kind: Secret
      jsonPointers:
        - /data/tls.key
        - /data/tls.crt
    - group: "apps"
      kind: Deployment
      jsonPointers:
        - /spec/template/metadata/annotations/checksum~1admission-webhook.yaml
```
