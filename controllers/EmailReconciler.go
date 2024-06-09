package controllers

import (
    "context"

    "github.com/go-logr/logr"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"

    examplev1 "github.com/toni-marlop/email-operator/api/v1"
)

type EmailReconciler struct {
    client.Client
    Log    logr.Logger
    Scheme *runtime.Scheme
}

// Reconcile is part of the main Kubernetes reconciliation loop
// which aims to move the current state of the cluster closer to the desired state.
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("email", req.NamespacedName)

    // Fetch the Email resource
    var email examplev1.Email
    if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
        if errors.IsNotFound(err) {
            // The resource was deleted
            return ctrl.Result{}, nil
        }
        log.Error(err, "unable to fetch Email")
        return ctrl.Result{}, err
    }

    // Fetch the sender configuration
    var emailSenderConfig examplev1.EmailSenderConfig
    if err := r.Get(ctx, client.ObjectKey{
        Namespace: req.Namespace,
        Name:      email.Spec.SenderConfigRef,
    }, &emailSenderConfig); err != nil {
        log.Error(err, "unable to fetch EmailSenderConfig")
        return ctrl.Result{}, err
    }

    mailerClient := MailerClient{Client: r.Client}
    messageId, err := mailerClient.SendEmail(ctx, email, emailSenderConfig)
    if err != nil {
        log.Error(err, "failed to send email")
        email.Status.DeliveryStatus = "Failed"
        email.Status.Error = err.Error()
    } else {
        email.Status.DeliveryStatus = "Sent"
        email.Status.MessageId = messageId
    }

    if err := r.Status().Update(ctx, &email); err != nil {
        log.Error(err, "failed to update email status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.Email{}).
        Complete(r)
}

