package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	errSleepDuration = 10 * time.Minute
)

func main() {
	cfg := NewConfig()
	// creates the in-cluster config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func(cfg *Config, ctx context.Context) {
		for {
			password, err := getLoginPassword(ctx, cfg.AWS_ACCESS_KEY_ID, cfg.AWS_SECRET_ACCESS_KEY, cfg.AWS_REGION)
			if err != nil {
				log.Printf("[error] failed to get login password: %v, retrying in 10 minutes", err)
				time.Sleep(errSleepDuration)
				continue
			}

			// Delete the secret
			if err := clientset.CoreV1().Secrets(cfg.K8S_NAMESPACE).Delete(ctx, cfg.K8S_SECRET_NAME, metav1.DeleteOptions{}); err != nil {
				log.Printf("[error] failed to delete secret: %v", err)
				time.Sleep(errSleepDuration)
				continue
			}

			// Create the secret
			dockerConfigJson := `{
	        "auths": {
	            "` + cfg.DOCKER_SERVER + `": {
	                "username": "AWS",
	                "password": "` + password + `",
	                "email": "` + cfg.DOCKER_EMAIL + `",
	                "auth": "` + base64.StdEncoding.EncodeToString([]byte("AWS:"+password)) + `"
	            }
	        }
	    }`

			dockerConfigBase64 := base64.StdEncoding.EncodeToString([]byte(dockerConfigJson))
			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfg.K8S_SECRET_NAME,
				},
				Type: v1.SecretTypeDockerConfigJson,
				Data: map[string][]byte{
					v1.DockerConfigJsonKey: []byte(dockerConfigBase64),
				},
			}
			if _, err := clientset.CoreV1().Secrets(cfg.K8S_NAMESPACE).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
				log.Printf("[error] failed to create secret: %v", err)
				time.Sleep(errSleepDuration)
				continue
			}

			log.Println("Secret created successfully")

			// Sleep for 10 hours
			time.Sleep(10 * time.Hour)
		}
	}(cfg, ctx)

	<-ctx.Done()

	log.Println("Exiting")
}

// gets the login password for the ECR registry
// id and secret are the AWS access key and secret
// region is the AWS region
func getLoginPassword(ctx context.Context, id string, secret string, region string) (string, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(id, secret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	ecrClient := ecr.NewFromConfig(awsCfg)

	// get login password (auth token)
	base64EncodedPasswords, err := ecrClient.GetAuthorizationToken(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get authorization token: %v", err)
	}
	if len(base64EncodedPasswords.AuthorizationData) == 0 {
		return "", fmt.Errorf("no authorization data found")
	}

	// decode the base64 encoded password
	base64Password := base64EncodedPasswords.AuthorizationData[0].AuthorizationToken
	if base64Password == nil {
		return "", fmt.Errorf("no password found")
	}

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(*base64Password)))
	n, err := base64.StdEncoding.Decode(dst, []byte(*base64Password))
	if err != nil {
		return "", fmt.Errorf("failed to decoce base64 password: %v", err)
	}
	dst = dst[:n]

	// trim AWS: from the password
	password := strings.TrimPrefix(string(dst), "AWS:")
	password = strings.TrimSpace(password)

	return password, nil
}
