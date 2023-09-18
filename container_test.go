package di_test

import (
	"github.com/alexandremahdhaoui/di" //nolint:depguard
	. "github.com/onsi/ginkgo/v2"      //nolint:depguard
	. "github.com/onsi/gomega"         //nolint:depguard
	"strconv"                          //nolint:goimports,gofumpt,gci
)

var _ = Describe("Container", func() {
	var userDefinedContainerName string

	var assignedKey string
	var assignedValue string
	var assignedPointer *string

	var notAssignedKey string

	var userDefinedContainer di.Container
	var defaultContainer di.Container

	Describe("creating a new container", func() {
		It("should create a new container", func() {
			container := di.New("test-container")
			Expect(container).ShouldNot(BeNil())
		})

		It("should panic", func() {
			Expect(func() {
				_ = di.New("")
			}).To(Panic())
		})
	})

	Describe("setting a value to a container", func() {
		It("should successfully set a value", func() {
			err := di.Set(di.New("test"), di.NewValue[int]("test", new(int)))
			Expect(err).ShouldNot(HaveOccurred())

			err = di.Set(di.DefaultContainer, di.NewValue[int]("test", new(int)))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should fail setting a value to a built container", func() {
			c := di.New("test")
			c.Build()
			err := di.Set(c, di.NewValue[int]("test", new(int)))
			Expect(err).Should(HaveOccurred())
		})
	})

	BeforeEach(func() {
		userDefinedContainerName = "user-defined-container"

		assignedKey = "assigned-key"
		assignedValue = "assigned-value"
		assignedPointer = &assignedValue

		notAssignedKey = "not-assigned-key"

		userDefinedContainer = di.New(userDefinedContainerName)
		defaultContainer = di.DefaultContainer // Warn: this container will not be reset throughout the code

		err := di.Set(userDefinedContainer, di.NewValue(assignedKey, assignedPointer))
		if err != nil {
			panic(err)
		}
	})

	Describe("initializing values in an immature container", func() {
		It("should initialize values", func() {
			for _, c := range []di.Container{userDefinedContainer} {
				for i := 0; i < 10; i++ {
					di.MustWithOptions[int](c, strconv.Itoa(i), di.InitializeOption).MustSet(i)
				}
			}
		})
	})

	Describe("accessing an immature container", func() {
		Context("with Get", func() {
			It("should successfully return a value", func() {
				v, err := di.Get[string](userDefinedContainer, assignedKey) //nolint:nolintlint,varnamelen
				Expect(err).ShouldNot(HaveOccurred())
				Expect(v.MustPtr()).To(BeIdenticalTo(assignedPointer))
				Expect(v.MustValue()).To(BeIdenticalTo(assignedValue))

			})
		})
		Context("with Must", func() {
			It("should successfully return a value", func() {
				Expect(func() { di.Must[string](userDefinedContainer, assignedKey) }).ShouldNot(Panic())
				Expect(di.Must[string](userDefinedContainer, assignedKey).MustPtr()).To(BeIdenticalTo(assignedPointer))
				Expect(di.Must[string](userDefinedContainer, assignedKey).MustValue()).To(BeIdenticalTo(assignedValue))
			})
		})
		Context("with MustWithOptions", func() {
			It("should successfully return a value", func() {
				Expect(func() { di.MustWithOptions[string](userDefinedContainer, assignedKey) }).ShouldNot(Panic())
				Expect(di.MustWithOptions[string](userDefinedContainer, assignedKey).MustPtr()).
					To(BeIdenticalTo(assignedPointer))
				Expect(di.MustWithOptions[string](userDefinedContainer, assignedKey).MustValue()).
					To(BeIdenticalTo(assignedValue))
			})
		})
	})

	Describe("building a new container", func() {
		BeforeEach(func() {
			for _, c := range []di.Container{userDefinedContainer, defaultContainer} {
				for i := 0; i < 10; i++ {
					di.MustWithOptions[int](c, strconv.Itoa(i), di.InitializeOption).MustSet(i)
				}
			}
		})

		It("should build non empty container", func() {
			Expect(func() { userDefinedContainer.Build() }).NotTo(Panic())
			Expect(func() { defaultContainer.Build() }).NotTo(Panic())
		})
	})

	Describe("initializing values in a built container", func() {
		BeforeEach(func() {
			userDefinedContainer.Build()
		})

		Context("with Set", func() {
			It("should fail", func() {
				err := di.Set(userDefinedContainer, di.NewValue(assignedKey, assignedPointer))
				Expect(err).Should(HaveOccurred())

				err = di.Set(defaultContainer, di.NewValue(assignedKey, assignedPointer))
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("with MustWithOptions and InitializeOption", func() {
			It("should panic", func() {
				Expect(func() {
					di.MustWithOptions[int](userDefinedContainer, "should-panic", di.InitializeOption).MustSet(0)
				}).To(Panic())
			})
		})
	})

	Describe("accessing a built container's value", func() {
		BeforeEach(func() {
			userDefinedContainer.Build()
		})

		Context("with Get", func() {
			It("should successfully return a value", func() {
				v, err := di.Get[string](userDefinedContainer, assignedKey)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(v.MustValue()).Should(BeIdenticalTo(assignedValue))
				Expect(v.MustPtr()).Should(Equal(assignedPointer))
				// Should not be an identical pointer
				Expect(v.MustPtr()).ShouldNot(BeIdenticalTo(assignedPointer))
			})

			It("should fail to return a value", func() {
				_, err := di.Get[int](defaultContainer, notAssignedKey)
				Expect(err).Should(HaveOccurred())

				_, err = di.Get[int](userDefinedContainer, notAssignedKey)
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("with Must", func() {
			It("should successfully return a value", func() {
				v := di.Must[string](userDefinedContainer, assignedKey)
				Expect(v.MustValue()).Should(BeIdenticalTo(assignedValue))
				Expect(v.MustPtr()).Should(Equal(assignedPointer))
				Expect(v.MustPtr()).ShouldNot(BeIdenticalTo(assignedPointer))
			})

			It("should panic retrieving an unassigned value", func() {
				Expect(func() { di.Must[string](userDefinedContainer, notAssignedKey) }).To(Panic())
				Expect(func() { di.Must[string](defaultContainer, notAssignedKey) }).To(Panic())
			})
		})

		Context("with MustWithOptions", func() {
			It("should successfully return a value", func() {
				Expect(di.MustWithOptions[string](userDefinedContainer, assignedKey).MustValue()).
					To(Equal(assignedValue))
			})

			It("should panic on not assigned values", func() {
				Expect(func() { di.MustWithOptions[string](userDefinedContainer, notAssignedKey) }).To(Panic())
				Expect(func() { di.MustWithOptions[string](defaultContainer, notAssignedKey) }).To(Panic())
			})
		})
	})
})
