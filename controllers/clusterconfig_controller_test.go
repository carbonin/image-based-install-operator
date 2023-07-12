package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	cro "github.com/RHsyseng/cluster-relocation-operator/api/v1beta1"
	relocationv1alpha1 "github.com/carbonin/cluster-relocation-service/api/v1alpha1"
	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/filesystem"
	bmh_v1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Reconcile", func() {
	var (
		c               client.Client
		dataDir         string
		serverDir       string
		r               *ClusterConfigReconciler
		ctx             = context.Background()
		configName      = "test-config"
		configNamespace = "test-namespace"
	)

	BeforeEach(func() {
		c = fakeclient.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&relocationv1alpha1.ClusterConfig{}).
			Build()
		var err error
		dataDir, err = os.MkdirTemp("", "clusterconfig_controller_test_data")
		Expect(err).NotTo(HaveOccurred())
		serverDir, err = os.MkdirTemp("", "clusterconfig_controller_test_server")
		Expect(err).NotTo(HaveOccurred())

		r = &ClusterConfigReconciler{
			Client: c,
			Scheme: scheme.Scheme,
			Log:    logrus.New(),
			Options: &ClusterConfigReconcilerOptions{
				BaseURL:   "https://example.com/",
				DataDir:   dataDir,
				ServerDir: serverDir,
			},
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(dataDir)).To(Succeed())
		Expect(os.RemoveAll(serverDir)).To(Succeed())
	})

	createSecret := func(name string, data map[string][]byte) {
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "test-namespace",
			},
			Data: data,
		}
		Expect(c.Create(ctx, s)).To(Succeed())
	}

	validateISOSecretContent := func(fs filesystem.FileSystem, file string, data map[string][]byte) {
		f, err := fs.OpenFile(file, os.O_RDONLY)
		Expect(err).NotTo(HaveOccurred())

		content, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		secret := &corev1.Secret{}
		Expect(json.Unmarshal(content, secret)).To(Succeed())
		Expect(secret.Data).To(Equal(data))
	}

	It("creates an iso with the correct relocation content", func() {
		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
			Spec: relocationv1alpha1.ClusterConfigSpec{
				ClusterRelocationSpec: cro.ClusterRelocationSpec{
					Domain:  "thing.example.com",
					SSHKeys: []string{"ssh-rsa sshkeyhere foo@example.com"},
				},
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		key := types.NamespacedName{
			Namespace: configNamespace,
			Name:      configName,
		}
		res, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: key})
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		isoPath := filepath.Join(serverDir, configNamespace, fmt.Sprintf("%s.iso", configName))
		d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
		Expect(err).NotTo(HaveOccurred())
		fs, err := d.GetFilesystem(0)
		Expect(err).NotTo(HaveOccurred())
		f, err := fs.OpenFile("/cluster-relocation-spec.json", os.O_RDONLY)
		Expect(err).NotTo(HaveOccurred())

		content, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())
		relocationSpec := &cro.ClusterRelocationSpec{}
		Expect(json.Unmarshal(content, relocationSpec)).To(Succeed())
		Expect(*relocationSpec).To(Equal(config.Spec.ClusterRelocationSpec))
	})

	It("creates the referenced secrets", func() {
		apiCertData := map[string][]byte{"apicert": []byte("apicert")}
		ingressCertData := map[string][]byte{"ingresscert": []byte("ingresscert")}
		pullSecretData := map[string][]byte{"pullsecret": []byte("pullsecret")}
		createSecret("api-cert", apiCertData)
		createSecret("ingress-cert", ingressCertData)
		createSecret("pull-secret", pullSecretData)

		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
			Spec: relocationv1alpha1.ClusterConfigSpec{
				ClusterRelocationSpec: cro.ClusterRelocationSpec{
					APICertRef: &corev1.SecretReference{
						Name: "api-cert", Namespace: configNamespace,
					},
					IngressCertRef: &corev1.SecretReference{
						Name: "ingress-cert", Namespace: configNamespace,
					},
					PullSecretRef: &corev1.SecretReference{
						Name: "pull-secret", Namespace: configNamespace,
					},
				},
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		key := types.NamespacedName{
			Namespace: configNamespace,
			Name:      configName,
		}
		res, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: key})
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		isoPath := filepath.Join(serverDir, configNamespace, fmt.Sprintf("%s.iso", configName))
		d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
		Expect(err).NotTo(HaveOccurred())
		fs, err := d.GetFilesystem(0)
		Expect(err).NotTo(HaveOccurred())

		validateISOSecretContent(fs, "/api-cert-secret.json", apiCertData)
		validateISOSecretContent(fs, "/ingress-cert-secret.json", ingressCertData)
		validateISOSecretContent(fs, "/pull-secret-secret.json", pullSecretData)
	})

	It("sets the image url in status", func() {
		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		key := types.NamespacedName{
			Namespace: configNamespace,
			Name:      configName,
		}
		res, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: key})
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(c.Get(ctx, key, config)).To(Succeed())
		Expect(config.Status.ImageURL).To(Equal(fmt.Sprintf("https://example.com/images/%s/%s.iso", configNamespace, configName)))
	})

	It("configures a referenced BMH", func() {
		bmh := &bmh_v1alpha1.BareMetalHost{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bmh",
				Namespace: "test-bmh-namespace",
			},
		}
		Expect(c.Create(ctx, bmh)).To(Succeed())

		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
			Spec: relocationv1alpha1.ClusterConfigSpec{
				BareMetalHostRef: &relocationv1alpha1.BareMetalHostReference{
					Name:      bmh.Name,
					Namespace: bmh.Namespace,
				},
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: configNamespace,
				Name:      configName,
			},
		}
		res, err := r.Reconcile(ctx, req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		key := types.NamespacedName{
			Namespace: bmh.Namespace,
			Name:      bmh.Name,
		}
		Expect(c.Get(ctx, key, bmh)).To(Succeed())
		Expect(bmh.Spec.Image).NotTo(BeNil())
		Expect(bmh.Spec.Image.URL).To(Equal(fmt.Sprintf("https://example.com/images/%s/%s.iso", configNamespace, configName)))
		Expect(bmh.Spec.Image.DiskFormat).To(HaveValue(Equal("live-iso")))
		Expect(bmh.Spec.Online).To(BeTrue())
	})
})

var _ = Describe("mapBMHToCC", func() {
	var (
		c               client.Client
		r               *ClusterConfigReconciler
		ctx             = context.Background()
		configName      = "test-config"
		configNamespace = "test-namespace"
	)

	BeforeEach(func() {
		c = fakeclient.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&relocationv1alpha1.ClusterConfig{}).
			Build()

		r = &ClusterConfigReconciler{
			Client: c,
			Scheme: scheme.Scheme,
			Log:    logrus.New(),
		}
	})

	It("returns a request for the cluster config referencing the given BMH", func() {
		bmh := &bmh_v1alpha1.BareMetalHost{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bmh",
				Namespace: "test-bmh-namespace",
			},
		}
		Expect(c.Create(ctx, bmh)).To(Succeed())

		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
			Spec: relocationv1alpha1.ClusterConfigSpec{
				BareMetalHostRef: &relocationv1alpha1.BareMetalHostReference{
					Name:      bmh.Name,
					Namespace: bmh.Namespace,
				},
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		config = &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-config",
				Namespace: configNamespace,
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		requests := r.mapBMHToCC(ctx, bmh)
		Expect(len(requests)).To(Equal(1))
		Expect(requests[0].NamespacedName).To(Equal(types.NamespacedName{
			Name:      configName,
			Namespace: configNamespace,
		}))
	})

	It("returns an empty list when no cluster config matches", func() {
		bmh := &bmh_v1alpha1.BareMetalHost{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bmh",
				Namespace: "test-bmh-namespace",
			},
		}
		Expect(c.Create(ctx, bmh)).To(Succeed())

		config := &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: configNamespace,
			},
			Spec: relocationv1alpha1.ClusterConfigSpec{
				BareMetalHostRef: &relocationv1alpha1.BareMetalHostReference{
					Name:      "other-bmh",
					Namespace: bmh.Namespace,
				},
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())

		config = &relocationv1alpha1.ClusterConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-config",
				Namespace: configNamespace,
			},
		}
		Expect(c.Create(ctx, config)).To(Succeed())
		requests := r.mapBMHToCC(ctx, bmh)
		Expect(len(requests)).To(Equal(0))
	})
})