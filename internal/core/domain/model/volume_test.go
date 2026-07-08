package model

import (
	"delivery/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VolumeBeCorrectWhenParamAreCorrectOnCreate(t *testing.T) {
	expectedV := uint16(10)

	volume, err := NewVolume(expectedV)

	assert.NoError(t, err)
	assert.NotEmpty(t, volume)
	assert.Equal(t, expectedV, volume.volume)
}

func Test_VolumeReturnErrorWhenParamIsInCorrectOnCreated(t *testing.T) {
	v := uint16(0)
	expected := errs.NewValueMustBeGreaterThan("volume", v, 0)

	_, err := NewVolume(v)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}
