package git

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func Test_ToUpper(t *testing.T) {
	t.Parallel()
	// Test empty string early return optimization
	require.Equal(t, "", ToUpper(""))

	require.Equal(t, "/MY/NAME/IS/:PARAM/*", ToUpper("/my/name/is/:param/*"))

	// Test single character optimizations
	require.Equal(t, "A", ToUpper("a"))
	require.Equal(t, "Z", ToUpper("z"))
	require.Equal(t, "1", ToUpper("1")) // non-letter remains unchanged
	require.Equal(t, "!", ToUpper("!")) // special character remains unchanged
}