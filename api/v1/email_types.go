package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

var (
    GroupVersion = schema.GroupVersion{Group: "example.com", Version: "v1"}
    SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
    AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
    scheme.AddKnownTypes(GroupVersion,
        &Email{},
        &EmailList{},
        &EmailSenderConfig{},
        &EmailSenderConfigList{},
    )
    metav1.AddToGroupVersion(scheme, GroupVersion)
    return nil
}

// Email CRD
type EmailSpec struct {
    SenderConfigRef string `json:"senderConfigRef"`
    RecipientEmail  string `json:"recipientEmail"`
    Subject         string `json:"subject"`
    Body            string `json:"body"`
}

type EmailStatus struct {
    DeliveryStatus string `json:"deliveryStatus"`
    MessageID      string `json:"messageId"`
    Error          string `json:"error,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Email struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   EmailSpec   `json:"spec,omitempty"`
    Status EmailStatus `json:"status,omitempty"`
}

func (e *Email) DeepCopyObject() runtime.Object {
    return &Email{
        TypeMeta:   e.TypeMeta,
        ObjectMeta: *e.ObjectMeta.DeepCopy(),
        Spec:       e.Spec,
        Status:     e.Status,
    }
}

// +kubebuilder:object:root=true
type EmailList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Email `json:"items"`
}

func (el *EmailList) DeepCopyObject() runtime.Object {
    return &EmailList{
        TypeMeta: el.TypeMeta,
        ListMeta: *el.ListMeta.DeepCopy(),
        Items:    append([]Email(nil), el.Items...),
    }
}

// EmailSenderConfig CRD
type EmailSenderConfigSpec struct {
    ApiTokenSecretRef string `json:"apiTokenSecretRef"`
    SenderEmail       string `json:"senderEmail"`
    Provider          string `json:"provider"`
}

type EmailSenderConfigStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type EmailSenderConfig struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
    Status EmailSenderConfigStatus `json:"status,omitempty"`
}

func (esc *EmailSenderConfig) DeepCopyObject() runtime.Object {
    return &EmailSenderConfig{
        TypeMeta:   esc.TypeMeta,
        ObjectMeta: *esc.ObjectMeta.DeepCopy(),
        Spec:       esc.Spec,
        Status:     esc.Status,
    }
}

// +kubebuilder:object:root=true
type EmailSenderConfigList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []EmailSenderConfig `json:"items"`
}

func (escl *EmailSenderConfigList) DeepCopyObject() runtime.Object {
    return &EmailSenderConfigList{
        TypeMeta: escl.TypeMeta,
        ListMeta: *escl.ListMeta.DeepCopy(),
        Items:    append([]EmailSenderConfig(nil), escl.Items...),
    }
}

func init() {
    SchemeBuilder.Register(addKnownTypes)
}

