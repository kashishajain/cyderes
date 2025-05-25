package main

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "encoding/json"
    "github.com/kashishajain/cyderes-app/transformer"
)

func TestTransform(t *testing.T) {
    orig := transformer.Original_data{
        UserID: 1,
        ID:     101,
        Title:  "Test Title",
        Body:   "Test Body",
    }
    input, err := json.Marshal(orig)
    assert.NoError(t, err)

    result, err := transformer.TransformData(input)

    assert.Equal(t, orig.UserID, result.UserID)
    assert.Equal(t, orig.ID, result.ID)
    assert.Equal(t, orig.Title, result.Title)
    assert.Equal(t, orig.Body, result.Body)
    assert.Equal(t, "placeholder_api", result.Source)

    now := time.Now().UTC()
    assert.WithinDuration(t, now, result.IngestedAt, time.Second*2, "ingested_at timestamp is not recent")
}
