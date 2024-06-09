# Email Operator

This repository contains a Kubernetes operator for managing custom resources to configure and send emails using transactional email providers like MailerSend and Mailgun.

## Custom Resource Definitions (CRDs)

### EmailSenderConfig

Defines the configuration for sending emails.

```yaml
apiVersion: example.com/v1
kind: EmailSenderConfig
metadata:
  name: my-email-config
  namespace: email-operator
spec:
  apiTokenSecretRef: my-email-secret
  senderEmail: sender@example.com
```

## Email
Defines an email to be sent.

```yaml
apiVersion: example.com/v1
kind: Email
metadata:
  name: test-email
  namespace: email-operator
spec:
  senderConfigRef: my-email-config
  recipientEmail: recipient@example.com
  subject: Test Email
  body: This is a test email.
status:
  deliveryStatus: ""
  messageId: ""
  error: ""
```
## EmailSenderConfig
### Prerequisites
- Kubernetes cluster (Minikube can be used for local testing)
- Kubectl installed and configured
- Docker installed and running
### Steps
1- Build the Docker image:
```bash
docker build -t email-operator:latest .
```
2- Load the image into Minikube:
```bash
minikube image load email-operator:latest
```
3- Deploy the Custom Resource Definitions (CRDs):
```bash
kubectl apply -f config/crd/bases/example.com_emails.yaml
kubectl apply -f config/crd/bases/example.com_emailsenderconfigs.yaml
```
4- Deploy the operator:
```bash
kubectl apply -f config/manager/manager.yaml
```
5- Create the necessary resources:
- Create a Secret containing the MailerSend API token:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-email-secret
  namespace: email-operator
data:
  apiToken: <base64-encoded-api-token>

```
- Apply the Secret:
```bash
kubectl apply -f secret.yaml
```
- Create an EmailSenderConfig resource:
```bash
kubectl apply -f emailsenderconfig.yaml
```
- Create an Email resource:
```bash
kubectl apply -f email.yaml
```
## Usage
Monitor the logs of the operator to verify that emails are being processed and sent:
```bash
kubectl logs deployment/email-operator -n email-operator
```
Check the status of an Email resource to see the delivery status, message ID, and any errors:
```bash
kubectl get email test-email -n email-operator -o yaml
```





