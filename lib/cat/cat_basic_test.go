package cat_test

import (
	"github.com/yarbelk/slimbox/lib/cat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"strings"
	"testing"
)

var c *cat.CatOptions = &cat.CatOptions{}

var _ = Describe("Basic", func() {
	var (
		C *cat.CatOptions = &cat.CatOptions{}
		inputValue string
		inputReader *strings.Reader
		outputStream *bytes.Buffer
	)

	BeforeEach(func() {
		outputStream = new(bytes.Buffer)
	})

	Context("Empty File", func() {
		BeforeEach(func() {
			inputValue = ""
			inputReader = strings.NewReader(inputValue)
		})

		It("return an empty response", func() {
			C.Cat(inputReader, outputStream)
			Expect(outputStream.String()).To(Equal(""))
		})
	})

	Context("One line file", func() {
		BeforeEach(func() {
			inputValue = "Hello World!\n"
			inputReader = strings.NewReader(inputValue)
		})

		It("return the same as the input", func() {
			C.Cat(inputReader, outputStream)
			Expect(outputStream.String()).To(Equal(inputValue))
		})
	})

	Context("One line file without trailing newline", func() {
		BeforeEach(func() {
			inputValue = "Hello World!"
			inputReader = strings.NewReader(inputValue)
		})

		It("the input with a trailing newline", func() {
			C.Cat(inputReader, outputStream)
			Expect(outputStream.String()).To(Equal("Hello World!\n"))
		})
	})

	Context("with multiple lines", func() {
		BeforeEach(func() {
			inputValue = "Hello World!\n\nHello Gophers!\n"
			inputReader = strings.NewReader(inputValue)
		})

		It("Should be output on multiple lines the same data", func() {
			C.Cat(inputReader, outputStream)
			Expect(outputStream.String()).To(Equal(inputValue))
		})
	})

	Context("when calling cat twice", func() {
		var (
			firstInputValue string
			secondInputValue string
			inputReaderOne *strings.Reader
			inputReaderTwo *strings.Reader
		)

		Context("with files ending in newlines", func() {
			BeforeEach(func() {
				firstInputValue = "Hello World!\n"
				secondInputValue = "Hello Gophers!\n"
				inputReaderOne = strings.NewReader(firstInputValue)
				inputReaderTwo = strings.NewReader(secondInputValue)
				
			})

			It("Concatinates both inputs", func() {
				C.Cat(inputReaderOne, outputStream)
				C.Cat(inputReaderTwo, outputStream)

				Expect(outputStream.String()).To(Equal("Hello World!\nHello Gophers!\n"))
			})

		Context("wiht one file without a newline", func() {
			BeforeEach(func() {
				firstInputValue = "Hello World!"
				secondInputValue = "Hello Gophers!\n"
				inputReaderOne = strings.NewReader(firstInputValue)
				inputReaderTwo = strings.NewReader(secondInputValue)
				
			})

			It("Concatinates both inputs", func() {
				C.Cat(inputReaderOne, outputStream)
				C.Cat(inputReaderTwo, outputStream)

				Expect(outputStream.String()).To(Equal("Hello World!\nHello Gophers!\n"))
			})
		})
		})
	})
})
