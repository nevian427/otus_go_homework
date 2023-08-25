package main

import (
	"crypto/rand"
	mrand "math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testdataSize = 32 * 1024

func TestCopyRegular(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		err := os.Remove("testfile_copy")
		require.NoErrorf(t, err, "Error remove copy file: %s", err)
	})

	// make "golden" data to compare with
	testdata := make([]byte, testdataSize)
	copybuffer := make([]byte, testdataSize)

	if _, err := rand.Read(testdata); err != nil {
		t.Fatalf("Create golden data error: %s", err)
	}

	testfile, err := os.CreateTemp("", "otusgo_hw07_")
	if err != nil {
		t.Fatalf("Create test file error: %s", err)
	}
	defer os.Remove(testfile.Name()) // clean up

	if _, err := testfile.Write(testdata); err != nil {
		t.Fatalf("Error write test file: %s", err)
	}
	if err := testfile.Close(); err != nil {
		t.Fatalf("Error close test file: %s", err)
	}

	t.Run("zero offset and rand limit", func(t *testing.T) {
		limit := mrand.Int63n(testdataSize)

		err := Copy(testfile.Name(), "testfile_copy", 0, limit)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		require.Equal(t, limit, int64(n), "written bytes not equal red")
		require.Equal(t, testdata[:limit], copybuffer[:limit])
	})

	t.Run("rand offset and zero limit", func(t *testing.T) {
		offset := mrand.Int63n(testdataSize)

		err := Copy(testfile.Name(), "testfile_copy", offset, 0)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		size := testdataSize - offset

		require.Equal(t, size, int64(n), "written bytes not equal red")
		require.Equal(t, testdata[offset:], copybuffer[:size])
	})

	t.Run("rand offset and rand limit", func(t *testing.T) {
		limit := mrand.Int63n(testdataSize)
		offset := mrand.Int63n(testdataSize)

		err := Copy(testfile.Name(), "testfile_copy", offset, limit)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		size := 32*1024 - offset
		if size < limit {
			limit = size
		}
		require.Equal(t, limit, int64(n), "written bytes not equal red")
		require.Equal(t, testdata[offset:offset+limit], copybuffer[:limit])
	})

	t.Run("rand negative offset and rand limit", func(t *testing.T) {
		limit := mrand.Int63n(testdataSize)
		offset := -mrand.Int63n(testdataSize)

		err := Copy(testfile.Name(), "testfile_copy", offset, limit)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		if limit > -offset {
			limit = -offset
		}

		require.Equal(t, limit, int64(n), "written bytes not equal red")
		// размер минус отрицательное смещение даёт плюс
		require.Equal(t, testdata[testdataSize+offset:testdataSize+offset+limit], copybuffer[:limit])
	})

	t.Run("too big offset", func(t *testing.T) {
		// >testdataSize
		err := Copy(testfile.Name(), "testfile_copy", testdataSize+1024, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("negative limit", func(t *testing.T) {
		err := Copy(testfile.Name(), "testfile_copy", 0, -2)
		require.ErrorIs(t, err, ErrLimitNegative)
	})

	t.Run("default limit & offset", func(t *testing.T) {
		err := Copy(testfile.Name(), "testfile_copy", 0, 0)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		require.Equal(t, testdataSize, n, "written bytes not equal red")
		require.Equal(t, testdata, copybuffer)
	})
}

func TestCopyIrregular(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		err := os.Remove("testfile_copy_irr")
		require.NoErrorf(t, err, "Error remove copy file: %s", err)
	})

	t.Run("negative offset", func(t *testing.T) {
		err := Copy("/dev/urandom", "testfile_copy_irr", -2, 10)
		require.ErrorIs(t, err, ErrOffsetInvalid)
	})

	t.Run("zero limit", func(t *testing.T) {
		err := Copy("/dev/urandom", "testfile_copy_irr", 0, 0)
		require.ErrorIs(t, err, ErrLimitInvalid)
	})

	t.Run("rand offset and rand limit", func(t *testing.T) {
		copybuffer := make([]byte, testdataSize)

		limit := mrand.Int63n(testdataSize)
		offset := mrand.Int63n(testdataSize)

		err := Copy("/dev/urandom", "testfile_copy_irr", offset, limit)
		require.NoErrorf(t, err, "Copy func error: %s", err)

		testcopy, err := os.Open("testfile_copy_irr")
		require.NoErrorf(t, err, "Error open copy file: %s", err)
		defer testcopy.Close()

		n, err := testcopy.Read(copybuffer)
		require.NoErrorf(t, err, "Copy read error: %s", err)

		require.Equal(t, limit, int64(n), "written bytes not equal red")
	})
}
