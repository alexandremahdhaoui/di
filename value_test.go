package di_test

import (
	"github.com/alexandremahdhaoui/di" //nolint:depguard
	. "github.com/onsi/ginkgo/v2"      //nolint:depguard
	. "github.com/onsi/gomega"         //nolint:depguard
)

type (
	testStruct struct {
		a int
	}

	testConcreteIface struct {
		b string
	}
	testIface interface {
		test() string
	}
)

func (t *testConcreteIface) test() string {
	return t.b
}

var _ = Describe("Value", func() {
	Describe("NewValue", func() {
		It("should return a Value[T]", func() {})
		It("should panic", func() {})
	})

	// Describe helps organize our Specs.
	// Describe node's closure can contain any number of Setup nodes (e.g. BeforeEach, AfterEach, JustBeforeEach), and
	// 	Subject nodes (i.e. It).
	// Context and When nodes are aliases for Describe - use whichever gives your suite a better narrative flow.
	// 	It is idiomatic to Describe the behavior of an object or function and, within that Describe, outline a number of
	//	Contexts and Whens.
	Describe("Value[T]", func() {
		var v0Inner []bool
		// v1Inner is not referenced
		var v2Inner testStruct
		var v3Inner testIface

		var v0InnerPointer *[]bool
		var v2InnerPointer *testStruct
		var v3InnerPointer *testIface

		var value0 di.Value[[]bool]
		var value1 di.Value[string]
		var value2 di.Value[testStruct]
		var value3 di.Value[testIface]

		// We instantiate states of our Specs
		BeforeEach(func() {
			v0Inner = []bool{true, false}
			// v1Inner is not referenced
			v2Inner = testStruct{a: 2}
			v3Inner = testIface(&testConcreteIface{b: "value3"})

			v0InnerPointer = &v0Inner
			v2InnerPointer = &v2Inner
			v3InnerPointer = &v3Inner

			value0 = di.NewValue[[]bool]("value0", v0InnerPointer)
			value1 = di.NewValue[string]("value1", nil)
			value2 = di.NewValue[testStruct]("value2", v2InnerPointer)
			value3 = di.NewValue[testIface]("value3", v3InnerPointer)
		})

		// Context is an alias to Describe but is used to structure the Specs
		Context("Key", func() {
			// "It" nodes are subject nodes that contains the Spec code & assertions
			It("should return the exact same key", func() {
				Expect(value0.Key()).To(Equal("value0"))
			})
		})

		Context("Ptr", func() {
			It("should return the exact same pointer", func() {
				ptr0, err := value0.Ptr()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ptr0).To(BeIdenticalTo(v0InnerPointer))

				ptr2, err := value2.Ptr()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ptr2).To(BeIdenticalTo(v2InnerPointer))

				ptr3, err := value3.Ptr()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ptr3).To(BeIdenticalTo(v3InnerPointer))
			})

			It("should return a pointer to a new string", func() {
				ptr1, err := value1.Ptr()
				Expect(err).ShouldNot(HaveOccurred())
				newStringPtr := new(string)
				Expect(ptr1).To(Equal(newStringPtr))
			})

		})

		Context("Set", func() {
			It("should have occurred", func() {
				newItem := []bool{false}
				err := value0.Set(newItem)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("Value", func() {
			It("should return the same value", func() {
				v0, err := value0.Value()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(v0).To(Equal(v0Inner))

				v1, err := value1.Value()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(v1).To(Equal(*new(string)))
			})
		})
	})

})
