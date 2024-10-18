// Handling State Retrieval

// Assuming you have outputs stored in S3 as described earlier
func (e *executor) retrieveTerraformOutputs(ctx context.Context, options Options) (map[string]interface{}, error) {
    outputsPath := "/tmp/outputs.json"
    err := downloadFromS3(ctx, "your-outputs-bucket", "outputs.json", outputsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to download outputs.json: %v", err)
    }

    outputsData, err := os.ReadFile(outputsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read outputs.json: %v", err)
    }

    var outputs map[string]interface{}
    err = json.Unmarshal(outputsData, &outputs)
    if err != nil {
        return nil, fmt.Errorf("failed to parse outputs.json: %v", err)
    }

    return outputs, nil
}

// Example: Reading the state file from a mounted path
func (e *executor) readTerraformState(stateFilePath string) (map[string]interface{}, error) {
    stateData, err := os.ReadFile(stateFilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read state file: %v", err)
    }

    var state map[string]interface{}
    err = json.Unmarshal(stateData, &state)
    if err != nil {
        return nil, fmt.Errorf("failed to parse state file: %v", err)
    }

    return state, nil
}
