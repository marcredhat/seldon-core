package controllers

import (
	"encoding/base64"
	"encoding/json"
	machinelearningv1alpha2 "github.com/seldonio/seldon-core/operator/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sort"
	"strings"
)

// Get the Namespace from the SeldonDeployment. Returns "default" if none found.
func getNamespace(deployment *machinelearningv1alpha2.SeldonDeployment) string {
	if len(deployment.ObjectMeta.Namespace) > 0 {
		return deployment.ObjectMeta.Namespace
	} else {
		return "default"
	}
}

// Translte the PredictorSpec p in to base64 encoded JSON.
func getEngineVarJson(p *machinelearningv1alpha2.PredictorSpec) (string, error) {
	str, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(str), nil
}

// Get an environment variable given by key or return the fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get an annotation from the Seldon Deployment given by annotationKey or return the fallback.
func getAnnotation(mlDep *machinelearningv1alpha2.SeldonDeployment, annotationKey string, fallback string) string {
	if annotation, hasAnnotation := mlDep.Spec.Annotations[annotationKey]; hasAnnotation {
		return annotation
	} else {
		return fallback
	}
}

//get annotations that start with seldon.io/engine
func getEngineEnvAnnotations(mlDep *machinelearningv1alpha2.SeldonDeployment) []corev1.EnvVar {

	envVars := make([]corev1.EnvVar, 0)
	var keys []string
	for k, _ := range mlDep.Spec.Annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		//prefix indicates engine annotation but "seldon.io/engine-separate-pod" isn't an env one
		if strings.HasPrefix(k, "seldon.io/engine-") && k != machinelearningv1alpha2.ANNOTATION_SEPARATE_ENGINE {
			name := strings.TrimPrefix(k, "seldon.io/engine-")
			var replacer = strings.NewReplacer("-", "_")
			name = replacer.Replace(name)
			name = strings.ToUpper(name)
			envVars = append(envVars, corev1.EnvVar{Name: name, Value: mlDep.Spec.Annotations[k]})
		}
	}
	return envVars
}
