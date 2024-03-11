package profile

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ProfileSpec struct {
	Type         string       `json:"type"`
	LocationSpec LocationSpec `json:"locationSpec"`
	Infra        Infra        `json:"infra"`
}

type Infra struct {
	Type       string     `json:"type"`
	OpenStack  OpenStack  `json:"openStack"`
	Azure      Azure      `json:"azure"`
	Portworx   Portworx   `json:"portworx"`
	Vsphere    Vsphere    `json:"vsphere"`
	Credential Credential `json:"credential"`
}

type Vsphere struct {
	ServerAddress string `json:"serverAddress"`
}

type Portworx struct {
	Namespace   string `json:"namespace"`
	ServiceName string `json:"serviceName"`
}

type Azure struct {
	ADEndpoint              string `json:"ADEndpoint"`
	ADResource              string `json:"ADResource"`
	CloudEnv                string `json:"cloudEnv"`
	CredentialType          string `json:"credentialType"`
	ResourceGroup           string `json:"resourceGroup"`
	ResourceManagerEndpoint string `json:"resourceManagerEndpoint"`
	SubscriptionID          string `json:"subscriptionID"`
	UseDefaultMSI           string `json:"useDefaultMSI"`
}

type OpenStack struct {
	KeystoneEndpoint string `json:"keystoneEndpoint"`
}

type LocationSpec struct {
	Credential    Credential `json:"credential"`
	Location      Location   `json:"location"`
	InfraPortable bool       `jon:"infraPortable"`
}

type Location struct {
	LocationType string      `json:"locationType"`
	ObjectStore  ObjectStore `json:"objectStore"`
	FileStore    FileStore   `json:"fileStore"`
	Vbr          Vbr         `json:"vbr"`
}

type Vbr struct {
	ServerAddress string `json:"serverAddress"`
	ServerPort    string `json:"serverPort"`
	RepoName      string `json:"repoName"`
	RepoId        string `json:"repoId"`
	SkipSSLVerify bool   `json:"skipSSLVerify"`
}

type FileStore struct {
	ClaimName string `json:"claimName"`
	Path      string `json:"path"`
}

type ObjectStore struct {
	ObjectStoreType  string `json:"objectStoreType"`
	Endpoint         string `json:"endpoint"`
	SkipSSLVerify    bool   `json:"skipSSLVerify"`
	Name             string `json:"name"`
	Region           string `json:"region"`
	Path             string `json:"path"`
	PathType         string `json:"pathType"`
	ProtectionPeriod string `json:"protectionPeriod"`
}

type Credential struct {
	SecretType string `json:"secretType"`
	Secret     Secret `json:"secret"`
}

type Secret struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
}

type ProfileStatus struct {
	Validation string   `json:"validation"`
	Hash       int64    `json:"hash"`
	Error      []string `json:"error"`
}

type Profile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileSpec   `json:"spec"`
	Status ProfileStatus `json:"status"`
}

type ProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Profile `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Profile) DeepCopyInto(out *Profile) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Status = ProfileStatus{
		Validation: in.Status.Validation,
		Hash:       in.Status.Hash,
		Error:      in.Status.Error,
	}
	out.Spec = ProfileSpec{
		Type: in.Spec.Type,
		LocationSpec: LocationSpec{
			Credential: Credential{
				SecretType: in.Spec.LocationSpec.Credential.SecretType,
				Secret: Secret{
					ApiVersion: in.Spec.LocationSpec.Credential.Secret.ApiVersion,
					Kind:       in.Spec.LocationSpec.Credential.Secret.Kind,
					Name:       in.Spec.LocationSpec.Credential.Secret.Name,
					Namespace:  in.Spec.LocationSpec.Credential.Secret.Namespace,
				},
			},
			Location: Location{
				LocationType: in.Spec.LocationSpec.Location.LocationType,
				ObjectStore: ObjectStore{
					ObjectStoreType:  in.Spec.LocationSpec.Location.ObjectStore.ObjectStoreType,
					Endpoint:         in.Spec.LocationSpec.Location.ObjectStore.Endpoint,
					SkipSSLVerify:    in.Spec.LocationSpec.Location.ObjectStore.SkipSSLVerify,
					Name:             in.Spec.LocationSpec.Location.ObjectStore.Name,
					Region:           in.Spec.LocationSpec.Location.ObjectStore.Region,
					Path:             in.Spec.LocationSpec.Location.ObjectStore.Path,
					PathType:         in.Spec.LocationSpec.Location.ObjectStore.PathType,
					ProtectionPeriod: in.Spec.LocationSpec.Location.ObjectStore.ProtectionPeriod,
				},
				FileStore: FileStore{
					ClaimName: in.Spec.LocationSpec.Location.FileStore.ClaimName,
					Path:      in.Spec.LocationSpec.Location.FileStore.Path,
				},
				Vbr: Vbr{
					ServerAddress: in.Spec.LocationSpec.Location.Vbr.ServerAddress,
					ServerPort:    in.Spec.LocationSpec.Location.Vbr.ServerPort,
					RepoName:      in.Spec.LocationSpec.Location.Vbr.RepoName,
					RepoId:        in.Spec.LocationSpec.Location.Vbr.RepoId,
					SkipSSLVerify: in.Spec.LocationSpec.Location.Vbr.SkipSSLVerify,
				},
			},
			InfraPortable: in.Spec.LocationSpec.InfraPortable,
		},
		Infra: Infra{
			Type: in.Spec.Infra.Type,
			OpenStack: OpenStack{
				KeystoneEndpoint: in.Spec.Infra.OpenStack.KeystoneEndpoint,
			},
			Azure: Azure{
				ADEndpoint:              in.Spec.Infra.Azure.ADEndpoint,
				ADResource:              in.Spec.Infra.Azure.ADResource,
				CloudEnv:                in.Spec.Infra.Azure.CloudEnv,
				CredentialType:          in.Spec.Infra.Azure.CredentialType,
				ResourceGroup:           in.Spec.Infra.Azure.ResourceGroup,
				ResourceManagerEndpoint: in.Spec.Infra.Azure.ResourceManagerEndpoint,
				SubscriptionID:          in.Spec.Infra.Azure.SubscriptionID,
				UseDefaultMSI:           in.Spec.Infra.Azure.UseDefaultMSI,
			},
			Portworx: Portworx{
				Namespace:   in.Spec.Infra.Portworx.Namespace,
				ServiceName: in.Spec.Infra.Portworx.ServiceName,
			},
			Vsphere: Vsphere{
				ServerAddress: in.Spec.Infra.Vsphere.ServerAddress,
			},
			Credential: Credential{
				SecretType: in.Spec.Infra.Credential.SecretType,
				Secret: Secret{
					ApiVersion: in.Spec.Infra.Credential.Secret.ApiVersion,
					Kind:       in.Spec.Infra.Credential.Secret.Kind,
					Name:       in.Spec.Infra.Credential.Secret.Name,
					Namespace:  in.Spec.Infra.Credential.Secret.Namespace,
				},
			},
		},
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Profile) DeepCopyObject() runtime.Object {
	out := Profile{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ProfileList) DeepCopyObject() runtime.Object {
	out := ProfileList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Profile, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
