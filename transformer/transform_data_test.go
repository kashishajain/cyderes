package transformer

import (
    "fmt"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "encoding/json"
)

func TestTransform(t *testing.T) {
    orig := []Original_data{
        {
        UserID: 1,
        ID:     101,
        Title:  "Test Title",
        Body:   "Test Body",
        },
    }
    input, err := json.Marshal(orig)
    assert.NoError(t, err)

    //result := []transformer.Transformed_data{}
    result, err := TransformData(input)
    fmt.Println(result)
    assert.Equal(t, orig[0].UserID, result[0].UserID)
    assert.Equal(t, orig[0].ID, result[0].ID)
    assert.Equal(t, orig[0].Title, result[0].Title)
    assert.Equal(t, orig[0].Body, result[0].Body)
    assert.Equal(t, "placeholder_api", result[0].Source)

    now := time.Now().UTC()
    ingestedAt, err := time.Parse(time.RFC3339, result[0].IngestedAt)
    assert.NoError(t, err)
    assert.WithinDuration(t, now, ingestedAt, time.Second*2, "ingested_at timestamp is not recent")
}
