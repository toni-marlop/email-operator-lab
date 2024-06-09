package controllers

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/go-resty/resty/v2"
    examplev1 "github.com/toni-marlop/email-operator/api/v1"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

const mailerSendAPIURL = "https://api.mailersend.com/v1/email"

type MailerClient struct {
    client.Client
}

type MailerSendEmailRequest struct {
    From     EmailAddress `json:"from"`
    To       []EmailAddress `json:"to"`
    Subject  string       `json:"subject"`
    Text     string       `json:"text"`
}

type EmailAddress struct {
    Email string `json:"email"`
    Name  string `json:"name,omitempty"`
}

func (m *MailerClient) SendEmail(ctx context.Context, email examplev1.Email, senderConfig examplev1.EmailSenderConfig) (string, error) {
    var apiTokenSecret corev1.Secret
    if err := m.Get(ctx, client.ObjectKey{Namespace: email.Namespace, Name: senderConfig.Spec.APITokenSecretRef}, &apiTokenSecret); err != nil {
        return "", fmt.Errorf("failed to get API token secret: %w", err)
    }

    apiToken := string(apiTokenSecret.Data["apiToken"])

    client := resty.New()

    emailRequest := MailerSendEmailRequest{
        From: EmailAddress{Email: senderConfig.Spec.SenderEmail},
        To:   []EmailAddress{{Email: email.Spec.RecipientEmail}},
        Subject: email.Spec.Subject,
        Text: email.Spec.Body,
    }

    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+apiToken).
        SetHeader("Content-Type", "application/json").
        SetBody(emailRequest).
        Post(mailerSendAPIURL)

    if err != nil {
        return "", fmt.Errorf("failed to send email: %w", err)
    }

    if resp.IsError() {
        return "", fmt.Errorf("email send request failed: %s", resp.Status())
    }

    var result map[string]interface{}
    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    messageId, ok := result["messageId"].(string)
    if !ok {
        return "", fmt.Errorf("message ID not found in response")
    }

    return messageId, nil
}

