package deployments

import (
	"context"

	appsv1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	restclient "k8s.io/client-go/rest"
)

// // KubeEnvironment - How to connect to Kubernetes Cluster for Deployment
// // TODO: Still not sure what info needs to be here... username/password
// //       or token auth?  support both probably
// // NOTE: How to handle incluster
// type KubeEnvironment struct {
// 	Name      string       `json:"name"`
// 	Namespace string       `json:"namespace"`
// 	Env       string       `json:"env"`
// 	UserAuth  KubeUserAuth `json:"user"`
// 	Cluster   KubeCluster  `json:"cluster"`
// }

// // KubeUserAuth - User Auth Data Pulled from DB
// type KubeUserAuth struct {
// 	Name                  string `json:"name"`
// 	ClientCertificate     string `json:"clientcert"`
// 	ClientCertificateData []byte `json:"clientcert-data"`
// 	ClientKey             string `json:"clientkey"`
// 	ClientKeyData         []byte `json:"clientkey-data"`
// 	Token                 string `json:"token"`
// 	Username              string `json:"username"`
// 	Password              string `json:"password"`
// }

// // KubeCluster - Kubernetes Cluster Info, pulled from DB
// type KubeCluster struct {
// 	Name                     string `json:"name"`
// 	ServerUri                string `json:"server_uri"`
// 	CertificateAuthority     string `json:"certauth"`
// 	CertificateAuthorityData []byte `json:"certauth-data"`
// 	InSecure                 bool   `json:"insecure"`
// }

// GetCurrentDeployments - List Deployments...
func GetCurrentDeployments(ctx context.Context, config *restclient.Config) (appsv1beta1.DeploymentList, error) {
	//var config k8s.Config
	//config.Clusters = []ke.Cluster
	//config.UserAuth = []ke.UserAuth
	//config.Context

	//config, err := NewClient(ke)
	//if err != nil {
	//	return nil, fmt.Errorf("Error Creating K8s Client: %v", err)
	//}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		
	}

	var deployments appsv1beta1.DeploymentList
	if err := config.List(ctx, config.Namespace, &deployments); err != nil {
		return nil, err
	}
	return &deployments, nil
}

// // CreateNameSpace - Create Project Deployment Namespace
// func CreateNameSpace(ke KubeEnvironment) error {
// 	client, err := newClient(ke)
// 	if err != nil {
// 		return fmt.Errorf("Error Creating K8s Client: %v", err)
// 	}

// 	ns := corev1.Namespace{
// 		Metadata: &metav1.ObjectMeta{
// 			Name: k8s.String(ke.Namespace),
// 		},
// 	}
// 	err = client.Create(context.Background(), &ns)
// 	return nil
// }

// func CreateDeployment(ke KubeEnvironemnt, deployment interface) error {
// 	client, err := newClient(ke)
// 	if err != nil {
// 		return fmt.Errorf("Error Create K8s Client: %v", err)
// 	}

// 	deployment := appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:
// 		}
// 	}
// }
