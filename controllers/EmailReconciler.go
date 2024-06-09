package controllers

import (
    "context"
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/go-logr/logr"
    "github.com/mailgun/mailgun-go/v4"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/types"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    examplev1 "github.com/toni-marlop/email-operator/api/v1"
)

type EmailReconciler struct {
    client.Client
    Log logr.Logger
}

func (r *EmailReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
    ctx := context.Background()
    log := r.Log.WithValues("email", req.NamespacedName)

    // 1. Get the Email resource that triggered the reconciliation event.
    var email examplev1.Email
    if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
        log.Error(err, "unable to fetch Email")
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 2. Get the EmailSenderConfig resource referenced by the Email.
    var config examplev1.EmailSenderConfig
    if err := r.Get(ctx, types.NamespacedName{Name: email.Spec.SenderConfigRef, Namespace: req.Namespace}, &config); err != nil {
        log.Error(err, "unable to fetch EmailSenderConfig")
        return ctrl.Result{}, err
    }

    // 3. Get the secret that contains the API token.
    secret := &corev1.Secret{}
    if err := r.Get(ctx, types.NamespacedName{Name: config.Spec.ApiTokenSecretRef, Namespace: req.Namespace}, secret); err != nil {
        log.Error(err, "unable to fetch Secret")
        return ctrl.Result{}, err
    }

    apiToken := string(secret.Data["apiToken"])

    var err error
    // 4. Select the provider based on the configuration and call the appropriate function to send the email.
    switch config.Spec.Provider {
    case "mailersend":
        err = sendMailWithMailerSend(apiToken, config.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body)
    case "mailgun":
        err = sendMailWithMailgun(apiToken, config.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body)
    default:
        err = fmt.Errorf("unsupported email provider: %s", config.Spec.Provider)
    }

    // 5. Update the Email resource status based on the email send result.
    if err != nil {
        log.Error(err, "failed to send email")
        email.Status.DeliveryStatus = "Failed"
        email.Status.Error = err.Error()
    } else {
        email.Status.DeliveryStatus = "Sent"
        email.Status.MessageId = "message-id-placeholder" // You should get the real messageId from the provider's response
    }

    // 6. Update the Email resource status in Kubernetes.
    if err := r.Status().Update(ctx, &email); err != nil {
        log.Error(err, "unable to update Email status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

// Function to send email using MailerSend
func sendMailWithMailerSend(apiToken, senderEmail, recipientEmail, subject, body string) error {
    url := "https://api.mailersend.com/v1/email"
    payload := map[string]interface{}{
        "from": map[string]string{"email": senderEmail},
        "to": []map[string]string{{"email": recipientEmail}},
        "subject": subject,
        "text": body,
    }
    jsonData, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+apiToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
    }
    return nil
}

// Function to send email using Mailgun
func sendMailWithMailgun(apiToken, senderEmail, recipientEmail, subject, body string) error {
    domain := "your-mailgun-domain.com"
    m := mailgun.NewMailgun(domain, apiToken)
    message := m.NewMessage(senderEmail, subject, body, recipientEmail)

    _, _, err := m.Send(context.Background(), message)
    return err
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.Email{}).
        Complete(r)
}
