package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var inputBytes = []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore
et dolore magna aliqua. Donec ac odio tempor orci dapibus ultrices in iaculis. At auctor urna nunc id cursus metus 
aliquam eleifend mi. Enim ut tellus elementum sagittis vitae et leo. Eget aliquet nibh praesent tristique magna sit 
amet purus. Sagittis orci a scelerisque purus. Ac feugiat sed lectus vestibulum. Sed velit dignissim sodales ut. 
Odio ut enim blandit volutpat maecenas volutpat blandit aliquam. Nibh sit amet commodo nulla facilisi nullam vehicula 
ipsum a. Mattis aliquam faucibus purus in massa tempor nec feugiat. Habitant morbi tristique senectus et. Interdum 
velit laoreet id donec ultrices tincidunt.Accumsan sit amet nulla facilisi morbi tempus. Eget arcu dictum varius 
duis at consectetur. Morbi tristique senectus et netus et malesuada fames ac. Curabitur vitae nunc sed velit 
dignissim. Imperdiet dui accumsan sit amet nulla facilisi morbi tempus. Mollis nunc sed id semper risus in hendrerit 
gravida. At auctor urna nunc id cursus metus aliquam eleifend mi. Porttitor leo a diam sollicitudin tempor id eu. 
Volutpat ac tincidunt vitae semper quis. Dignissim cras tincidunt lobortis feugiat vivamus at augue eget. Lacus vel 
facilisis volutpat est velit. Consectetur adipiscing elit ut aliquam purus sit. Morbi enim nunc faucibus a pellentesque 
sit. Proin libero nunc consequat interdum varius sit amet mattis.
Ac tortor vitae purus faucibus ornare suspendisse sed. Faucibus scelerisque eleifend donec pretium vulputate sapien nec 
sagittis aliquam. Sed enim ut sem viverra. Semper eget duis at tellus at urna condimentum mattis pellentesque. Dignissim 
diam quis enim lobortis scelerisque fermentum. Maecenas accumsan lacus vel facilisis volutpat est velit egestas dui. 
Blandit massa enim nec dui nunc mattis. Tellus molestie nunc non blandit massa enim nec. Risus nullam eget felis eget 
nunc lobortis mattis aliquam. Nibh tortor id aliquet lectus proin. Nunc id cursus metus aliquam eleifend. Nulla facilisi 
morbi tempus iaculis. Quis varius quam quisque id diam vel quam. Maecenas sed enim ut sem. Tincidunt ornare massa eget 
egestas purus viverra. Nunc id cursus metus aliquam eleifend mi in nulla. Leo in vitae turpis massa sed elementum tempus 
egestas. Turpis cursus in hac habitasse platea dictumst quisque sagittis purus. Consectetur lorem donec massa sapien 
faucibus et molestie. Integer vitae justo eget magna fermentum iaculis. Turpis massa sed elementum tempus egestas sed 
sed. Quam quisque id diam vel. Aliquet porttitor lacus luctus accumsan tortor posuere ac ut. Morbi tristique senectus 
et netus et. Consectetur adipiscing elit pellentesque habitant morbi tristique senectus et. Ante metus dictum at tempor 
commodo ullamcorper a lacus. Bibendum est ultricies integer quis auctor elit sed vulputate mi. Ultricies leo integer 
malesuada nunc vel risus commodo viverra maecenas. Non consectetur a erat nam at. Adipiscing tristique risus nec 
feugiat in fermentum posuere urna. At augue eget arcu dictum varius duis at consectetur. Ipsum suspendisse ultrices 
gravida dictum fusce ut placerat orci. Arcu risus quis varius quam quisque id diam. Adipiscing diam donec adipiscing 
tristique risus nec feugiat in. Adipiscing vitae proin sagittis nisl rhoncus. Tincidunt augue interdum velit euismod 
in pellentesque massa. Habitant morbi tristique senectus et netus et malesuada fames. Vestibulum mattis ullamcorper 
velit sed. Sit amet facilisis magna etiam tempor orci eu lobortis. Donec ac odio tempor orci. Duis tristique 
sollicitudin nibh sit amet. Dictum non consectetur a erat. In nisl nisi scelerisque eu ultrices. Nisl purus in mollis 
nunc sed. Viverra nibh cras pulvinar mattis nunc sed blandit libero. Sit amet aliquam id diam maecenas ultricies mi 
eget. Feugiat in fermentum posuere urna nec tincidunt praesent semper. Et netus et malesuada fames ac turpis egestas sed.`)

var inputBytesLen = int64(len(inputBytes))

func TestCopy(t *testing.T) {
	inputPath := "input.txt"
	outputPath := "output.txt"
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return
	}
	defer func() {
		inputFile.Close()
		os.Remove(inputPath)
	}()

	_, err = inputFile.Write(inputBytes)
	if err != nil {
		return
	}

	t.Run("copy all", func(t *testing.T) {
		err = Copy(inputPath, outputPath, 0, 0)
		require.NoError(t, err)
		defer os.Remove(outputPath)

		outputFile, err := os.Open(outputPath)
		require.NoError(t, err)

		outputFileInfo, err := outputFile.Stat()
		require.NoError(t, err)

		require.Equal(t, inputBytesLen, outputFileInfo.Size())
	})

	t.Run("copy with offset", func(t *testing.T) {
		var tOffset int64 = 350

		err = Copy(inputPath, outputPath, tOffset, 0)
		require.NoError(t, err)
		defer os.Remove(outputPath)

		outputFile, err := os.Open(outputPath)
		require.NoError(t, err)

		outputFileInfo, err := outputFile.Stat()
		require.NoError(t, err)

		require.Equal(t, inputBytesLen-tOffset, outputFileInfo.Size())
	})

	t.Run("copy with limit", func(t *testing.T) {
		var tLimit int64 = 100

		err = Copy(inputPath, outputPath, 0, tLimit)
		require.NoError(t, err)
		defer os.Remove(outputPath)

		outputFile, err := os.Open(outputPath)
		require.NoError(t, err)

		outputFileInfo, err := outputFile.Stat()
		require.NoError(t, err)

		require.Equal(t, tLimit, outputFileInfo.Size())
	})

	t.Run("copy with offset and limit", func(t *testing.T) {
		var tOffset int64 = 350
		var tLimit int64 = 1000

		err = Copy(inputPath, outputPath, tOffset, tLimit)
		require.NoError(t, err)
		defer os.Remove(outputPath)

		outputFile, err := os.Open(outputPath)
		require.NoError(t, err)

		outputFileInfo, err := outputFile.Stat()
		require.NoError(t, err)

		require.Equal(t, tLimit, outputFileInfo.Size())
	})

	t.Run("offset more than file size", func(t *testing.T) {
		var tOffset int64 = 4000
		err = Copy(inputPath, outputPath, tOffset, 0)
		require.Errorf(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("limit more than file size", func(t *testing.T) {
		var tLimit int64 = 10000

		err = Copy(inputPath, outputPath, 0, tLimit)
		require.NoError(t, err)
		defer os.Remove(outputPath)

		outputFile, err := os.Open(outputPath)
		require.NoError(t, err)

		outputFileInfo, err := outputFile.Stat()
		require.NoError(t, err)

		require.Equal(t, inputBytesLen, outputFileInfo.Size())
	})

	t.Run("unsupported or empty file", func(t *testing.T) {
		err = Copy("/dev/random", outputPath, 0, 0)
		require.Errorf(t, err, ErrUnsupportedFile.Error())
	})
}
