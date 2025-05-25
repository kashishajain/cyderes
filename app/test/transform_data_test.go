package transformer_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "yourmodule/internal/transformer"
)

func TestTransform(t *testing.T) {
    input := transformer.APIPost{
        UserID: 1,
        ID:     101,
        Title:  "Test Title",
        Body:   "Test Body",
    }

    result := transformer.Transform(input)

    assert.Equal(t, input.UserID, result.UserID)
    assert.Equal(t, input.ID, result.ID)
    assert.Equal(t, input.Title, result.Title)
    assert.Equal(t, input.Body, result.Body)
    assert.Equal(t, "placeholder_api", result.Source)

    now := time.Now().UTC()
    assert.WithinDuration(t, now, result.IngestedAt, time.Second*2, "ingested_at timestamp is not recent")
}
