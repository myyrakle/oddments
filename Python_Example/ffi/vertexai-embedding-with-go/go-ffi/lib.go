package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"unsafe"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// EmbeddingResult represents the result of an embedding operation
type EmbeddingResult struct {
	TextEmbedding  []float32 `json:"text_embedding,omitempty"`
	ImageEmbedding []float32 `json:"image_embedding,omitempty"`
	VideoEmbedding []float32 `json:"video_embedding,omitempty"`
	Error          string    `json:"error,omitempty"`
}

// getTextEmbedding returns text embedding from Vertex AI
func getTextEmbedding(ctx context.Context, client *aiplatform.PredictionClient, projectID, location, text string) ([]float32, error) {
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/multimodalembedding@001", projectID, location)

	instanceValue, err := structpb.NewStruct(map[string]interface{}{
		"text": text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: []*structpb.Value{structpb.NewStructValue(instanceValue)},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to predict: %w", err)
	}

	if len(resp.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions returned")
	}

	prediction := resp.Predictions[0]
	predMap := prediction.GetStructValue().AsMap()

	embeddingInterface, ok := predMap["textEmbedding"]
	if !ok {
		return nil, fmt.Errorf("textEmbedding not found in response")
	}

	embeddingList, ok := embeddingInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("textEmbedding is not a list")
	}

	embeddings := make([]float32, len(embeddingList))
	for i, v := range embeddingList {
		if fv, ok := v.(float64); ok {
			embeddings[i] = float32(fv)
		}
	}

	return embeddings, nil
}

// getImageEmbedding returns image embedding from Vertex AI
func getImageEmbedding(ctx context.Context, client *aiplatform.PredictionClient, projectID, location, imageBytes string) ([]float32, error) {
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/multimodalembedding@001", projectID, location)

	instanceValue, err := structpb.NewStruct(map[string]interface{}{
		"image": map[string]interface{}{
			"bytesBase64Encoded": imageBytes,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: []*structpb.Value{structpb.NewStructValue(instanceValue)},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to predict: %w", err)
	}

	if len(resp.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions returned")
	}

	prediction := resp.Predictions[0]
	predMap := prediction.GetStructValue().AsMap()

	embeddingInterface, ok := predMap["imageEmbedding"]
	if !ok {
		return nil, fmt.Errorf("imageEmbedding not found in response")
	}

	embeddingList, ok := embeddingInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("imageEmbedding is not a list")
	}

	embeddings := make([]float32, len(embeddingList))
	for i, v := range embeddingList {
		if fv, ok := v.(float64); ok {
			embeddings[i] = float32(fv)
		}
	}

	return embeddings, nil
}

//export GetTextEmbedding
func GetTextEmbedding(projectID *C.char, location *C.char, text *C.char, credentialsJSON *C.char) *C.char {
	ctx := context.Background()

	goProjectID := C.GoString(projectID)
	goLocation := C.GoString(location)
	goText := C.GoString(text)
	goCredentials := C.GoString(credentialsJSON)

	var opts []option.ClientOption
	if goCredentials != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(goCredentials)))
	}

	client, err := aiplatform.NewPredictionClient(ctx, opts...)
	if err != nil {
		result := EmbeddingResult{Error: fmt.Sprintf("failed to create client: %v", err)}
		jsonBytes, _ := json.Marshal(result)
		return C.CString(string(jsonBytes))
	}
	defer client.Close()

	embeddings, err := getTextEmbedding(ctx, client, goProjectID, goLocation, goText)
	if err != nil {
		result := EmbeddingResult{Error: fmt.Sprintf("failed to get text embedding: %v", err)}
		jsonBytes, _ := json.Marshal(result)
		return C.CString(string(jsonBytes))
	}

	result := EmbeddingResult{TextEmbedding: embeddings}
	jsonBytes, _ := json.Marshal(result)
	return C.CString(string(jsonBytes))
}

//export GetImageEmbedding
func GetImageEmbedding(projectID *C.char, location *C.char, imageBase64 *C.char, credentialsJSON *C.char) *C.char {
	ctx := context.Background()

	goProjectID := C.GoString(projectID)
	goLocation := C.GoString(location)
	goImageBase64 := C.GoString(imageBase64)
	goCredentials := C.GoString(credentialsJSON)

	var opts []option.ClientOption
	if goCredentials != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(goCredentials)))
	}

	client, err := aiplatform.NewPredictionClient(ctx, opts...)
	if err != nil {
		result := EmbeddingResult{Error: fmt.Sprintf("failed to create client: %v", err)}
		jsonBytes, _ := json.Marshal(result)
		return C.CString(string(jsonBytes))
	}
	defer client.Close()

	embeddings, err := getImageEmbedding(ctx, client, goProjectID, goLocation, goImageBase64)
	if err != nil {
		result := EmbeddingResult{Error: fmt.Sprintf("failed to get image embedding: %v", err)}
		jsonBytes, _ := json.Marshal(result)
		return C.CString(string(jsonBytes))
	}

	result := EmbeddingResult{ImageEmbedding: embeddings}
	jsonBytes, _ := json.Marshal(result)
	return C.CString(string(jsonBytes))
}

//export FreeString
func FreeString(str unsafe.Pointer) {
	C.free(str)
}

func main() {}
